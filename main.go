package main

import (
	"fmt"
	"os"

	"github.com/dimitrijed93/dgtc/internal/utils"
	"github.com/dimitrijed93/dgtc/pkg/dgtc"
)

func main() {

	validateArgs()

	inPath := os.Args[1]
	outPath := os.Args[2]

	dgtc := dgtc.NewDgtc(inPath, outPath)
	dgtc.Start()

}

func validateArgs() {
	if os.Args == nil {
		fmt.Errorf("Params for In path and out path are required!")
	}
}

func validatePath(inPath string, outPath string) {
	if inPath == utils.EMPTY_STRING || outPath == utils.EMPTY_STRING {
		fmt.Errorf("In path and out path are required!")
	}
}
