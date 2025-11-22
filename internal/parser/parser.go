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

			switch nextRune{
			case ' ', 'E':
				exprBuilder.WriteRune(currRune)
				cf, err = t.NewWildCardFragment("*")
			case '/':
				exprBuilder.WriteRune(currRune)
				exprBuilder.WriteRune(nextRune)
				
				p.advance()
				nextRune = p.getCurrentToken()
				
				if unicode.IsDigit(nextRune) {
					num := p.readNumber()
					exprBuilder.WriteString(fmt.Sprintf("%d", num))
					cf, err = t.NewDivisorFragment(exprBuilder.String(), []uint8{num})
				}
			}
		case unicode.IsSpace(currRune):
			if unicode.IsSpace(p.peekNext()) {
				return nil, errors.New("malformed cron expression")
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

	num, _ := strconv.Atoi(p.input[p.currPos:c+1])
	p.currPos = c

	return uint8(num)
}
