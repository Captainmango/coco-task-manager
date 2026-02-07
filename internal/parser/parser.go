package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	ASTERISK        = '*'
	FORWARD_SLASH   = '/'
	COMMA           = ','
	DASH            = '-'
	END_OF_FRAGMENT = ' '
	END_OF_CRON     rune // inits to zero or string terminating char
)

type Parser struct {
	input              string
	currPos            uint8
	peekPos            uint8
	output             Cron
	exprBuilder        strings.Builder
	err                error
	inputLength        uint8
	nextFragmentIterFn func() (CronFragmentType, bool)
	stopFragmentIterFn func()
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

		nextFn, stopFn := p.output.ExpressionOrder()

		p.input = inputStr
		p.inputLength = uint8(len(inputStr))
		p.nextFragmentIterFn = nextFn
		p.stopFragmentIterFn = stopFn

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

func (p *Parser) Parse() (Cron, error) {
	for p.currPos < p.inputLength {
		var cf CronFragment

		currRune := p.getCurrentToken()
		switch {
		case currRune == ASTERISK:
			cf, p.err = p.handleWildCard()
		case unicode.IsDigit(currRune):
			cf, p.err = p.handleDigit()
		case unicode.IsSpace(currRune):
			peekToken := p.peekNext()
			if unicode.IsSpace(peekToken) {
				return Cron{}, ErrTooManySpaces(p.input, p.peekPos)
			}

			if unicode.IsLetter(peekToken) {
				return Cron{}, ErrMalformedCron(p.input, p.peekPos)
			}

			p.advance()
			continue
		default:
			p.err = ErrMalformedCron(p.input, p.peekPos)
		}

		if p.err != nil {
			return Cron{}, p.err
		}

		cf.FragmentType, p.err = p.getCurrentFragmentType()

		if p.err != nil {
			return Cron{}, p.err
		}

		p.output.Data = append(p.output.Data, cf)
		p.exprBuilder.Reset()
		p.advance()
	}

	return p.output, nil
}

func structureInputForParser(input any, shouldValidateLength bool) (string, error) {
	const EXPECTED_CRON_LENGTH = 5

	if out, ok := input.(string); ok {
		inputArr := strings.Split(out, " ")

		if shouldValidateLength && len(inputArr) != EXPECTED_CRON_LENGTH {
			return "", ErrInvalidInput(out)
		}

		return out, nil
	}

	if out, ok := input.([]string); ok {
		if shouldValidateLength && len(out) != EXPECTED_CRON_LENGTH {
			return "", ErrInvalidInput(out)
		}

		return strings.Join(out, " "), nil
	}

	// should never reach here
	return "", ErrInvalidInput(input)
}

func (p *Parser) handleDigit() (CronFragment, error) {
	var cf CronFragment
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
				cf, p.err = NewListFragment(p.exprBuilder.String(), nums)
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
				p.err = ErrMalformedCron(p.input, p.peekPos)
				done = true
			}
		}
	case DASH:
		p.exprBuilder.WriteRune(nextRune)
		p.advance()
		num2 := p.readNumber()
		nums = append(nums, num2)
		p.exprBuilder.WriteString(fmt.Sprintf("%d", num2))
		cf, _ = NewRangeFragment(p.exprBuilder.String(), nums)
	case END_OF_FRAGMENT, END_OF_CRON:
		cf, _ = NewSingleFragment(p.exprBuilder.String(), nums)
	default:
		err = ErrMalformedCron(p.input, p.peekPos)
	}

	return cf, err
}

func (p *Parser) handleWildCard() (CronFragment, error) {
	var cf CronFragment
	var err error

	currRune := p.getCurrentToken()
	p.advance()
	nextRune := p.getCurrentToken()

	switch nextRune {
	case END_OF_FRAGMENT, END_OF_CRON:
		p.exprBuilder.WriteRune(currRune)
		cf, err = NewWildCardFragment("*")
	case FORWARD_SLASH:
		p.exprBuilder.Write([]byte{byte(currRune), byte(nextRune)})

		p.advance()
		nextRune = p.getCurrentToken()

		if unicode.IsDigit(nextRune) {
			num := p.readNumber()
			p.exprBuilder.WriteString(fmt.Sprintf("%d", num))
			cf, err = NewDivisorFragment(p.exprBuilder.String(), []uint8{num})
		}
	default:
		err = ErrMalformedCron(p.input, p.peekPos)
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

func (p *Parser) getCurrentFragmentType() (CronFragmentType, error) {
	cType, ok := p.nextFragmentIterFn()

	if !ok {
		return "", ErrUnableToPullNextFragment()
	}

	return cType, nil
}
