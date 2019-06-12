package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/tidwall/gjson"
	// "github.com/asgavar/partia"
)

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

	dbName := gjson.Get(lines[0], "open.database").String()
	dbLogin := gjson.Get(lines[0], "open.login").String()
	dbPassword := gjson.Get(lines[0], "open.password").String()

	if isItFirstRun {
		createDbAndTables(dbName, dbLogin, dbPassword)
	}

	dbUri := fmt.Sprintf(
		"postgres://%s:%s@localhost/%s?sslmode=disable", dbLogin, dbPassword, dbName)

	db, err := sql.Open("postgres", dbUri)
	fmt.Println(err)
	fmt.Println(db)

	// 	for _, commandInvocation := range lines[1:] {
	// 	}
}

func createDbAndTables(db, login, password string) {
	dbAsInit, _ := sql.Open("postgres", "postgres://init:qwerty@localhost?sslmode=disable")

	createRoleSql := fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD '%s'", login, password)
	createDbSql := fmt.Sprintf("CREATE DATABASE %s", db)
	grantConnectSql := fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s", db, login)

	dbAsInit.Query(createRoleSql)
	dbAsInit.Query(createDbSql)
	dbAsInit.Query(grantConnectSql)

	// TODO create actual tables
}

func openDbConnection() {

}
