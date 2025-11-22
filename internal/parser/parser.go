package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	d "github.com/captainmango/coco-cron-parser/internal/data"
)

var (
	WILDCARD        = '*'
	DIVISOR         = '/'
	COMMA           = ','
	DASH            = '-'
	END_OF_FRAGMENT = ' '
	END_OF_CRON     rune // inits to zero. 0 is not a valid unicode char
)

type Parser struct {
	input       string
	currPos     uint8
	peekPos     uint8
	output      d.Cron
	exprBuilder strings.Builder
	err         error
	inputLength uint8
}

type ValidParserInput interface {
	string|[]string
}

func NewParser[T ValidParserInput](in T) (*Parser, error) {
	validInput, err := structureInputForParser[T](in)

	if err != nil {
		return nil, err
	}

	return &Parser{input: validInput, inputLength: uint8(len(validInput))}, nil
}

func (p *Parser) Parse() (d.Cron, error) {
	for p.currPos < p.inputLength {
		var cf d.CronFragment

		currRune := rune(p.input[p.currPos])
		switch {
		case currRune == WILDCARD:
			cf, p.err = p.handleWildCard()
		case unicode.IsDigit(currRune):
			cf, p.err = p.handleDigit()
		case unicode.IsSpace(currRune):
			peekToken := p.peekNext()
			if unicode.IsSpace(peekToken) {
				return nil, p.tooManySpacesErr()
			}

			if unicode.IsLetter(peekToken) {
				return nil, p.malformedCronErr()
			}

			p.advance()
			continue
		default:
			p.err = p.malformedCronErr()
		}

		if p.err != nil {
			return nil, p.err
		}

		p.output = append(p.output, cf)
		p.exprBuilder.Reset()
		p.advance()
	}

	return p.output, nil
}

func structureInputForParser[T ValidParserInput](input any) (string, error) {
	if out, ok := input.(string); ok {
		return out, nil
	}

	if out, ok := input.([]string); ok {
		return strings.Join(out, " "), nil
	}

	return "", fmt.Errorf("invalid input %v", input)
}

func (p *Parser) handleDigit() (d.CronFragment, error) {
	var cf d.CronFragment
	var err error

	num := p.readNumber()
	nums := []uint8{num}
	p.exprBuilder.WriteString(fmt.Sprintf("%d", num))

	p.advance()
	nextRune := p.getCurrentToken()
	done := false

	switch nextRune {
	case COMMA:
		p.exprBuilder.WriteRune(nextRune)
		p.advance()

		for !done {
			nextRune = p.getCurrentToken()
			switch {
			case nextRune == END_OF_FRAGMENT,
				nextRune == END_OF_CRON:
				cf, p.err = d.NewListFragment(p.exprBuilder.String(), nums)
				done = true
			case unicode.IsDigit(nextRune):
				num = p.readNumber()
				nums = append(nums, num)
				p.exprBuilder.WriteString(fmt.Sprintf("%d", num))
				p.advance()
			case nextRune == COMMA:
				p.exprBuilder.WriteRune(nextRune)
				p.advance()
			default:
				p.err = p.malformedCronErr()
				done = true
			}
		}
	case DASH:
		p.exprBuilder.WriteRune(nextRune)
		p.advance()
		num2 := p.readNumber()
		nums = append(nums, num2)
		p.exprBuilder.WriteString(fmt.Sprintf("%d", num2))
		cf, _ = d.NewRangeFragment(p.exprBuilder.String(), nums)
	case END_OF_FRAGMENT, END_OF_CRON:
		cf, _ = d.NewSingleFragment(p.exprBuilder.String(), nums)
	default:
		err = p.malformedCronErr()
	}

	return cf, err
}

func (p *Parser) handleWildCard() (d.CronFragment, error) {
	var cf d.CronFragment
	var err error

	currRune := p.getCurrentToken()
	p.advance()
	nextRune := p.getCurrentToken()

	switch nextRune {
	case END_OF_FRAGMENT, END_OF_CRON:
		p.exprBuilder.WriteRune(currRune)
		cf, err = d.NewWildCardFragment("*")
	case DIVISOR:
		p.exprBuilder.Write([]byte{byte(currRune), byte(nextRune)})

		p.advance()
		nextRune = p.getCurrentToken()

		if unicode.IsDigit(nextRune) {
			num := p.readNumber()
			p.exprBuilder.WriteString(fmt.Sprintf("%d", num))
			cf, err = d.NewDivisorFragment(p.exprBuilder.String(), []uint8{num})
		}
	default:
		err = p.malformedCronErr()
	}

	return cf, err
}

func (p *Parser) advance() {
	p.currPos += 1
	p.peekPos = p.currPos
}

func (p *Parser) getCurrentToken() rune {
	if p.currPos < p.inputLength {
		return rune(p.input[p.currPos])
	}

	return END_OF_CRON
}

func (p *Parser) peekNext() rune {
	if p.peekPos+1 < p.inputLength {
		return rune(p.input[p.peekPos+1])
	}

	return END_OF_CRON
}

func (p *Parser) readNumber() uint8 {
	for unicode.IsDigit(p.peekNext()) {
		p.peekPos += 1
	}

	num, err := strconv.Atoi(p.input[p.currPos : p.peekPos+1])
	if err != nil {
		p.err = err

		return 0
	}

	p.currPos = p.peekPos

	return uint8(num)
}

func (p *Parser) malformedCronErr() error {
	return fmt.Errorf("malformed cron expression: '%s' invalid character after postion %d", p.input, p.peekPos)
}

func (p *Parser) tooManySpacesErr() error {
	return fmt.Errorf("malformed cron expression: '%s' too many spaces at position %d", p.input, p.peekPos)
}
