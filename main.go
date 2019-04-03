package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mpchadwick/dbanon/src"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	requested := flag.String("config", "", "Configuration to use. magento2 is included out-of-box. Alternately, supply path to file")
	flag.Parse()

	config, err := dbanon.NewConfig(*requested)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	provider := dbanon.NewProvider()
	processor := dbanon.NewLineProcessor(config, provider)
	reader := bufio.NewReader(os.Stdin)

	// sqlparser can be noisy
	// https://github.com/xwb1989/sqlparser/blob/120387863bf27d04bc07db8015110a6e96d0146c/ast.go#L52
	// We don't want to hear about it
	log.SetOutput(ioutil.Discard)

	for {
		text, err := reader.ReadString('\n')
		fmt.Print(processor.ProcessLine(text))

		if err != nil {
			break
		}
	}
}