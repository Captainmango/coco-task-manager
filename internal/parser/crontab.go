package parser

import (
	"strings"

	"github.com/captainmango/coco-cron-parser/internal/data"
	"github.com/google/uuid"
)

type CrontabEntry struct {
	ID uuid.UUID
	Cron data.Cron
	Cmd string
}

func NewCrontabEntryFromString(input string) (CrontabEntry, error) {
	var err error
	var ctbE CrontabEntry
	parts := strings.Split(input, " root ")
	
	if len(parts) != 2 {
		return ctbE, invalidCronTabEntry(input)
	}

	cronPart := parts[0]

	p, err := NewParser(
		WithInput(cronPart, true),
	)

	if err != nil {
		return ctbE, err
	}

	// create unmarshaltext so we can just read it from the string?
	cron, err := p.Parse()
	if err != nil {
		return ctbE, err
	}

	moreParts := strings.Split(parts[1], " # ")
	if len(moreParts) != 2 {
		return ctbE, invalidCronTabEntry(input)
	}

	uuID, err := uuid.FromBytes([]byte(moreParts[1]))
	if err != nil {
		return ctbE, err
	}

	ctbE.Cron = cron
	ctbE.ID = uuID
	ctbE.Cmd = moreParts[1]

	return ctbE, nil
}
