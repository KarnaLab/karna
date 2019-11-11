package deploy

import (
	"errors"

	"github.com/karnalab/karna/core"
)

const (
	fileName = "karna.yml"
)

func getTargetFunction(config *core.KarnaConfigFile, target, alias *string) (function *core.KarnaFunction, err error) {
	for i, f := range config.Functions {
		if f.Name == *target {
			function = &config.Functions[i]
		}
	}

	if function == nil {
		err = errors.New("No function match the provided key")
		return
	}

	if ok := function.Aliases[*alias]; len(ok) == 0 {
		err = errors.New("Alias do not match any existing alias")
	}

	return
}
