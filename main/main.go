package main

import (
	"bufio"
	"io/ioutil"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/tidwall/gjson"
	"github.com/asgavar/partia"
)

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

	dbName := gjson.Get(lines[0], "open.database").String()
	dbLogin := gjson.Get(lines[0], "open.login").String()
	dbPassword := gjson.Get(lines[0], "open.password").String()

	if isItFirstRun {
		partia.CreateDbAndRole(dbName, dbLogin, dbPassword)
	}

	partia.OpenConnection(dbName, dbLogin, dbPassword)

	if isItFirstRun {
		sql, _ := ioutil.ReadFile("../setup-db.sql")
		_, err := partia.DB.Query(string(sql))
		fmt.Println(err)
	}

	for _, commandInvocation := range lines[1:] {
		cmd_output := partia.Dispatch(commandInvocation)
		fmt.Println(partia.RenderOutputJson(cmd_output))
	}
}
