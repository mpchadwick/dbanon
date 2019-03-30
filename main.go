package main

import "bufio"
import "fmt"
import "os"

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		fmt.Print(text)

		if err != nil {
			break
		}
	}
}