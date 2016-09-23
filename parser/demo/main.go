package main

import (
	"bufio"
	"fmt"
	"gaql/parser"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	for {
		line, err := in.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("ReadBytes: %s", err)
		}

		sline := strings.Trim(string(line), "\n")
		fmt.Printf("exp = %s\n", sline)
		parser.ExpParse(parser.NewMExpLexer(sline))
	}
}
