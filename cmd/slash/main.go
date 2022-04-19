package main

import (
	"encoding/hex"
	"flag"
	"io/ioutil"
	"os"

	"github.com/c0mm4nd/slashvm"
	"github.com/c0mm4nd/slashvm/configs"
)

var (
	styleFlag   string
	abiFileFlag string
	binFileFlag string
)

func init() {
	flag.StringVar(&styleFlag, "style", "ethereum", "the style of the compiler")
	flag.StringVar(&abiFileFlag, "abi file", "", "")
	flag.StringVar(&binFileFlag, "bin file", "", "")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	binFile, err := os.Open(binFileFlag)
	check(err)

	binInHex, err := ioutil.ReadAll(binFile)
	check(err)

	code, err := hex.DecodeString(string(binInHex))
	check(err)

	vm := slashvm.NewRuntime(code, nil, &slashvm.Context{}, configs.IstanbulConfig)
	handler := &SimpleHandler{}
	vm.Run(handler)
}
