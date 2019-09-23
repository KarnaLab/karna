package core

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

func LogSuccessMessage(message string) {
	fmt.Printf("%s %s\n",
		aurora.Bold(aurora.Green("> Karna:")),

		aurora.Bold(aurora.Blue(message)))
}

func LogErrorMessage(message string) {
	fmt.Printf("%s %s\n",
		aurora.Bold(aurora.Green("> Karna:")),

		aurora.Bold(aurora.Red(message)))
}
