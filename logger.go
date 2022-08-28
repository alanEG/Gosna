package main

import (
	"fmt"
	"os"
	"github.com/fatih/color"
)

func logger(typeLog, Flag, message, reason interface{}, errorRange int) {
	flagNoColor := flagNoColor

	if flagNoColor {
		color.NoColor = true // disables colorized output
	}

	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	switch typeLog {
	case "success":
		fmt.Fprintf(os.Stderr, "%s %s\n", yellow(Flag), message)
	case "error":
		if reason != "" && reason != nil {
			fmt.Fprintf(os.Stderr, "%s %s\n", red("[Error]"), message)
			fmt.Fprintf(os.Stderr, "Reason: %s\n", reason)
		} else {
			fmt.Fprintf(os.Stderr, "%s %s\n", red("[Error]"), message)
		}

		if errorRange == 1 {
			os.Exit(1)
		}
	}
}
