package main

import (
	"fmt"
	"karna/cmd"
	"karna/core"
)

func main() {
	cmd.Execute()
	fmt.Println(core.Lambda)
}
