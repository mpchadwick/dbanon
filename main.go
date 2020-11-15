package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/blang/semver"
	"github.com/sirupsen/logrus"
	"github.com/mpchadwick/dbanon/src"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"io/ioutil"
	"log"
	"runtime/pprof"
	"os"
)

var version string

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
	ver := flag.Bool("version", false, "Get current version")
	silent := flag.Bool("silent", false, "Disable all logging")
	logFile := flag.String("log-file", "", "File to write logs to")
	logLevel := flag.String("log-level", "", "Specify desired log level")
	profile := flag.Bool("profile", false, "Generate a profile")

	flag.Parse()

	if *profile {
		profF, _ := os.Create("dbanon.prof")
		pprof.StartCPUProfile(profF)
		defer pprof.StopCPUProfile()
	}

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	if *update {
		if err := selfUpdate(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	dbanonLogger := dbanon.GetLogger()
	if !*silent {
		f := "dbanon.log"
		if *logFile != "" {
			f = *logFile
		}
		file, _ := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		dbanonLogger.SetOutput(file)
	} else {
		dbanonLogger.SetOutput(ioutil.Discard)
	}

	if *logLevel != "" {
		level, err := logrus.ParseLevel(*logLevel)
		if err != nil {
			dbanonLogger.Error(err)
		} else {
			dbanonLogger.SetLevel(level)
		}
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
	mode := "anonymize"
	if len(args) > 0 && args[0] == "map-eav" {
		mode = "map-eav"
	}

	provider := dbanon.NewProvider()
	eav := dbanon.NewEav(config)
	processor := dbanon.NewLineProcessor(mode, config, provider, eav)


	for {
		text, err := reader.ReadString('\n')
		result := processor.ProcessLine(text)
		if mode == "anonymize" {
			fmt.Print(result)
		}

		if err != nil {
			break
		}
	}

	if mode == "map-eav" {
		out, _ := eav.Config.String()
		fmt.Print(string(out))
		os.Exit(0)
	}
}