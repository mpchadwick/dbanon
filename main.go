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

	// sqlparser can be noisy
	// https://github.com/xwb1989/sqlparser/blob/120387863bf27d04bc07db8015110a6e96d0146c/ast.go#L52
	// We don't want to hear about it
	log.SetOutput(ioutil.Discard)
	reader := bufio.NewReader(os.Stdin)


	args := flag.Args()
	if len(args) > 0 && args[0] == "map-eav" {
		eav := dbanon.NewEav(config)

		for {
			text, err := reader.ReadString('\n')
			eav.ProcessLine(text)

			if err != nil {
				break
			}
		}

		out, _ := eav.Config.String()
		fmt.Print(string(out))
		os.Exit(0)
	}

	provider := dbanon.NewProvider()
	processor := dbanon.NewLineProcessor(config, provider)

	for {
		text, err := reader.ReadString('\n')
		fmt.Print(processor.ProcessLine(text))

		if err != nil {
			break
		}
	}
}