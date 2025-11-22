package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	t "github.com/captainmango/coco-cron-parser/internal/cron"
)

type Parser struct {
	input   string
	currPos uint8
	peekPos uint8
}

func NewParser(in string) *Parser {
	return &Parser{input: in}
}

func (p *Parser) Parse() (t.Cron, error) {
	var output t.Cron
	var err error

	inputLength := len(p.input)
	var exprBuilder strings.Builder

	for p.currPos < uint8(inputLength) {
		var cf t.CronFragment

		currRune := rune(p.input[p.currPos])
		switch {
		case currRune == '*':
			p.advance()
			nextRune := p.getCurrentToken()

			switch nextRune {
			case ' ', 'E':
				exprBuilder.WriteRune(currRune)
				cf, err = t.NewWildCardFragment("*")
			case '/':
				exprBuilder.Write([]byte{byte(currRune), byte(nextRune)})

				p.advance()
				nextRune = p.getCurrentToken()

				if unicode.IsDigit(nextRune) {
					num := p.readNumber()
					exprBuilder.WriteString(fmt.Sprintf("%d", num))
					cf, err = t.NewDivisorFragment(exprBuilder.String(), []uint8{num})
				}
			}
		case unicode.IsDigit(currRune):
			num := p.readNumber()
			nums := []uint8{num}
			exprBuilder.WriteString(fmt.Sprintf("%d", num))

			p.advance()
			nextRune := p.getCurrentToken()
			done := false
			
			switch nextRune {
			case ',':
				exprBuilder.WriteRune(nextRune)
				p.advance()

				for !done {
					nextRune = p.getCurrentToken()
					switch {
					case nextRune == ' ',
						nextRune == 'E':
						cf, _ = t.NewListFragment(exprBuilder.String(), nums)
						done = true
					case unicode.IsDigit(nextRune):
						num = p.readNumber()
						nums = append(nums, num)
						exprBuilder.WriteString(fmt.Sprintf("%d", num))
						p.advance()
					case nextRune == ',':
						exprBuilder.WriteRune(nextRune)
						p.advance()
					default:
						err = fmt.Errorf("malformed cron expression: '%s' invalid character after postion %d", p.input, p.peekPos)
						done = true
					}
				}
			case '-':
				exprBuilder.WriteRune(nextRune)
				p.advance()
				num2 := p.readNumber()
				nums = append(nums, num2)
				exprBuilder.WriteString(fmt.Sprintf("%d", num2))
				cf, _ = t.NewRangeFragment(exprBuilder.String(), nums)
			case ' ', 'E':
				cf, _ = t.NewSingleFragment(exprBuilder.String(), nums)
			default:
				err = fmt.Errorf("malformed cron expression: '%s' invalid character after postion %d", p.input, p.peekPos)
			}

		case unicode.IsSpace(currRune):
			peekToken := p.peekNext()
			if unicode.IsSpace(peekToken) {
				return nil, fmt.Errorf("malformed cron expression: '%s' too many spaces at position %d", p.input, p.peekPos)
			}

			if unicode.IsLetter(peekToken) {
				return nil, fmt.Errorf("malformed cron expression: '%s' invalid character at postion %d", p.input, p.peekPos)
			}

			p.advance()
			continue
		default:
			err = errors.New("oops")
		}

		p.advance()

		if err != nil {
			return nil, err
		}

		output = append(output, cf)
		exprBuilder.Reset()
	}

	return output, nil
}

func (p *Parser) advance() {
	p.currPos += 1
	p.peekPos = p.currPos
}

func (p *Parser) getCurrentToken() rune {
	if p.currPos < uint8(len(p.input)) {
		return rune(p.input[p.peekPos])
	}

	return rune('E') // Signifies the end of the expression
}

func (p *Parser) peekNext() rune {
	p.peekPos += 1
	if p.peekPos < uint8(len(p.input)) {
		return rune(p.input[p.peekPos])
	}

	return rune('E') // Signifies the end of the expression
}

func (p *Parser) readNumber() uint8 {
	c := p.peekPos

	for unicode.IsDigit(p.peekNext()) {
		c += 1
	}

	num, _ := strconv.Atoi(p.input[p.currPos : c+1])
	p.currPos = c

	return uint8(num)
}
