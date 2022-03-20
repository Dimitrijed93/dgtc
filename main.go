package main

import (
	"log"
	"os"

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
	if len(os.Args) != 3 {
		log.Fatal("Invalid args. Input file and output directory are required")
	}
}
