package core

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

//LogSuccessMessage => Print into CLI a success message.
func LogSuccessMessage(message string) {
	fmt.Printf("%s %s\n",
		aurora.Bold(aurora.Green("> Karna:")),

		aurora.Bold(aurora.Blue(message)))
}

//LogSuccessMessage => Print into CLI an error message.
func LogErrorMessage(message string) {
	fmt.Printf("%s %s\n",
		aurora.Bold(aurora.Green("> Karna:")),

		aurora.Bold(aurora.Red(message)))
}
