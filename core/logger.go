package core

import (
	"fmt"
	"log"

	"github.com/logrusorgru/aurora"
)

//LogSuccessMessage => Print into CLI a success message.
func LogSuccessMessage(message string) {
	fmt.Printf("%s %s\n",
		aurora.Bold(aurora.Green("> Karna:")),

		aurora.Bold(aurora.Blue(message)))
}

//LogErrorMessage => Print into CLI an error message.
func LogErrorMessage(message string) {
	log.Fatalf("%s %s\n",
		aurora.Bold(aurora.Green("> Karna:")),

		aurora.Bold(aurora.Red(message)))
}
