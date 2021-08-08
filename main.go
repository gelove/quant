package main

import (
	"quant/cmd"
	_ "quant/pkg/orm"
)

func main() {
	cmd.Execute()
}
