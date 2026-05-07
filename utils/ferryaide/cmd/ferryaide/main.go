package main

import (
	"context"
	"fmt"
	"os"

	"ferryaide/internal/app"
	"ferryaide/internal/cli"
)

func main() {
	code := run(os.Args[1:])
	os.Exit(code)
}

func run(args []string) int {
	options, err := cli.Parse(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	options.Input = os.Stdin
	options.Output = os.Stdout
	result, err := app.Run(context.Background(), options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if result.Message != "" {
		fmt.Fprintln(os.Stdout, result.Message)
	}
	return 0
}
