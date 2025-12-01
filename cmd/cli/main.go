package main

import (
	"fmt"

	"github.com/captainmango/coco-cron-parser/internal/commands"
)

func main() {
	fmt.Println(commands.Registry.Commands...)
	// cronArgs := strings.Join(os.Args[1:], " ")

	// p, err := parser.NewParser(
	// 	parser.WithInput(cronArgs, true),
	// )

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// cron, _ := p.Parse()
	// cron.PrintingMode = data.RAW_EXPRESSION
	// fmt.Printf("%s\n", cron)
}
