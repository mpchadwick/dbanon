package main

import "bufio"
import "fmt"
import "github.com/mpchadwick/dbanon/src"
import "os"

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		fmt.Print(anomymizer.Anonymize(text))

		if err != nil {
			break
		}
	}
}