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
	trackerType := os.Args[3]

	dgtc := dgtc.NewDgtc(inPath, outPath, trackerType)
	dgtc.Start()

}

func validateArgs() {
	if len(os.Args) != 4 {
		log.Fatal("Invalid args. Input file and output directory are required")
	}
}
