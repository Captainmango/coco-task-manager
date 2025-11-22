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
	END_OF_CRON     rune // inits to zero or string terminating char
)

type Parser struct {
	input               string
	currPos             uint8
	peekPos             uint8
	output              d.Cron
	exprBuilder         strings.Builder
	err                 error
	inputLength         uint8
	currentFragmentIdx  uint8
	currentFragmentType d.CronFragmentType
}

type ValidParserInput interface {
	string | []string
}

type ParserOption func(p *Parser) error

func WithInput[T ValidParserInput](input T, shouldValidateLength bool) ParserOption {
	return func(p *Parser) error {
		var inputStr string
		var err error

		inputStr, err = structureInputForParser(input, shouldValidateLength)

		if err != nil {
			return err
		}

		p.input = inputStr
		p.inputLength = uint8(len(inputStr))
		p.currentFragmentIdx = 0
		p.currentFragmentType = p.output.ExpressionOrder()[p.currentFragmentIdx]

		return nil
	}
}

func NewParser(opts ...ParserOption) (*Parser, error) {
	p := Parser{}

	for _, opt := range opts {
		err := opt(&p)

		if err != nil {
			return nil, err
		}
	}

	return &p, nil
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

		cf.FragmentType = p.getCurrentFragmentType()
		p.output = append(p.output, cf)
		p.exprBuilder.Reset()
		p.setNextFragmentType()
		p.advance()
	}

	return p.output, nil
}

func structureInputForParser(input any, shouldValidateLength bool) (string, error) {
	const EXPECTED_CRON_LENGTH = 5

	if out, ok := input.(string); ok {
		inputArr := strings.Split(out, " ")

		if shouldValidateLength && len(inputArr) != EXPECTED_CRON_LENGTH {
			return "", fmt.Errorf("%s is not a valid input", out)
		}

		return out, nil
	}

	if out, ok := input.([]string); ok {
		if shouldValidateLength && len(out) != EXPECTED_CRON_LENGTH {
			return "", fmt.Errorf("%s is not a valid input", out)
		}

		return strings.Join(out, " "), nil
	}

	// should never reach here
	return "", fmt.Errorf("%s is not a valid input", input)
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

func (p *Parser) getCurrentFragmentType() d.CronFragmentType {
	return p.currentFragmentType
}

func (p *Parser) setNextFragmentType() {
	p.currentFragmentIdx += 1
	order := p.output.ExpressionOrder()
	
	if p.currentFragmentIdx >= uint8(len(order)) {
		p.currentFragmentIdx = uint8(len(order)) - 1
	}

	p.currentFragmentType = p.output.ExpressionOrder()[p.currentFragmentIdx]
}

func (p *Parser) malformedCronErr() error {
	return fmt.Errorf("malformed cron expression: '%s' invalid character after postion %d", p.input, p.peekPos)
}

func (p *Parser) tooManySpacesErr() error {
	return fmt.Errorf("malformed cron expression: '%s' too many spaces at position %d", p.input, p.peekPos)
}
