package roomapi

import (
	"github.com/debuconnor/dbcore"
)

var db dbcore.Connection

func checkDb() {
	err := db.ConnectMysql()
	if err != nil {
		Error(err)
	}
	defer db.DisconnectMysql()
}

func initDb() {
	db = dbcore.NewDb()
	db.SetConnectionFromGcpSecret(DB_SECRET_VERSION)
	checkDb()
	db.ConnectMysql()
}
