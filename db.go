package partia

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

func CreateDbAndRole(db, login, password string) {
	dbAsInit, _ := sql.Open("postgres", "postgres://init:qwerty@localhost?sslmode=disable")

	createRoleSql := fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD '%s'", login, password)
	createDbSql := fmt.Sprintf("CREATE DATABASE %s", db)
	grantConnectSql := fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s", db, login)

	dbAsInit.Query(createRoleSql)
	dbAsInit.Query(createDbSql)
	dbAsInit.Query(grantConnectSql)
}

func OpenConnection(db, login, password string) {
	dbUri := fmt.Sprintf(
		"postgres://%s:%s@localhost/%s?sslmode=disable", login, password, db)

	DB, _ = sql.Open("postgres", dbUri)
}
