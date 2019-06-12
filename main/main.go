package main

import "bufio"
import "encoding/json"
import "fmt"
import "os"

// import "github.com/asgavar/partia"

func main() {
	args := os.Args[1:]
	var isItFirstRun bool
	var lines []string

	if len(args) > 0 && args[0] == "--init" {
		isItFirstRun = true
	} else {
		isItFirstRun = false
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Println(isItFirstRun)

	for _, commandInvocation := range lines[1:] {
		var buf map[string]interface{}
		json.Unmarshal([]byte(commandInvocation), &buf)
		fmt.Println(buf)
	}
}
