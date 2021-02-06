package deploy

import (
	"log"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

type KarnaLogger struct{}

func (*KarnaLogger) Log(message string) {
	log.SetOutput(colorable.NewColorableStdout())
	log.Printf("%s %s\n",
		aurora.Bold(aurora.Blue("> Karna:")),

		aurora.Bold(aurora.Green("Success: "+message)))
}

func (*KarnaLogger) Error(message string) {
	log.SetOutput(colorable.NewColorableStdout())
	log.Printf("%s %s\n",
		aurora.Bold(aurora.Blue("> Karna:")),

		aurora.Bold(aurora.Red("Error: "+message)))
}
