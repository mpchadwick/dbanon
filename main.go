package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/blang/semver"
	"github.com/mpchadwick/dbanon/src"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"io/ioutil"
	"log"
	"os"
)

const version = "0.2.2"

const slug = "mpchadwick/dbanon"

func selfUpdate() error {
	previous := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(previous, slug)
	if err != nil {
		return err
	}

	if previous.Equals(latest.Version) {
		fmt.Println("Current binary is the latest version", version)
	} else {
		fmt.Println("Update successfully done to version", latest.Version)
		fmt.Println("Release note:\n", latest.ReleaseNotes)
	}
	return nil
}

func main() {
	requested := flag.String("config", "", "Configuration to use. magento2 is included out-of-box. Alternately, supply path to file")
	update := flag.Bool("update", false, "Auto update dbanon to the newest version")

	flag.Parse()

	if *update {
		if err := selfUpdate(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

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