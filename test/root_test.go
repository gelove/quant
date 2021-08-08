package test

import (
	"quant/cmd"
	"testing"
)

func TestInitConfig(t *testing.T) {
	t.Log("TestPath")
	cmd.InitConfig("/Users/allen/Projects/Go/mod/quant/config/config.yml")
}
