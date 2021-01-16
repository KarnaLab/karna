package deploy

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

type KarnaLogger struct{}

func (*KarnaLogger) Log(message string) {
	fmt.Printf("%s %s\n",
		aurora.Bold(aurora.Green("> Karna:")),

		aurora.Bold(aurora.Blue(message)))
}

func (*KarnaLogger) Error(message string) {
	fmt.Printf("%s %s\n",
		aurora.Bold(aurora.Green("> Karna:")),

		aurora.Bold(aurora.Red(message)))
}
