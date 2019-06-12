package main

import (
	"bufio"
	"database/sql"
	"io/ioutil"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/tidwall/gjson"
	"github.com/asgavar/partia"
)

var DB *sql.DB

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
		createDbAndRole(dbName, dbLogin, dbPassword)
	}

	dbUri := fmt.Sprintf(
		"postgres://%s:%s@localhost/%s?sslmode=disable", dbLogin, dbPassword, dbName)

	DB, _ = sql.Open("postgres", dbUri)

	if isItFirstRun {
		sql, _ := ioutil.ReadFile("../setup-db.sql")
		_, err := DB.Query(string(sql))
		fmt.Println(err)
	}

	for _, commandInvocation := range lines[1:] {
		partia.Dispatch(commandInvocation)
	}
}

func createDbAndRole(db, login, password string) {
	dbAsInit, _ := sql.Open("postgres", "postgres://init:qwerty@localhost?sslmode=disable")

	createRoleSql := fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD '%s'", login, password)
	createDbSql := fmt.Sprintf("CREATE DATABASE %s", db)
	grantConnectSql := fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s", db, login)

	dbAsInit.Query(createRoleSql)
	dbAsInit.Query(createDbSql)
	dbAsInit.Query(grantConnectSql)
}
