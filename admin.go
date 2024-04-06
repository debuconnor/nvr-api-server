package roomapi

import (
	"github.com/debuconnor/dbcore"
)

func NewAdmin(admin Admin) User {
	return &admin
}

func (user *Admin) Get() {
	dml := dbcore.NewDml()
	dml.SelectAll()
	dml.From(SCHEMA_ADMIN)
	dml.Where("", COLUMN_ID, dbcore.EQUAL, itoa(user.Id))
	queryResult := dml.Execute(db.GetDb())

	user.Id = atoi(queryResult[0][COLUMN_ID])
	user.UserId = queryResult[0][COLUMN_USER_ID]
	user.Password = queryResult[0][COLUMN_PASSWORD]
}

func (user *Admin) Save() {
	dml := dbcore.NewDml()
	dml.Insert()
	dml.Into(SCHEMA_ADMIN)
	dml.Value(COLUMN_USER_ID, user.UserId)
	dml.Value(COLUMN_PASSWORD, user.Password)
	dml.Execute(db.GetDb())

	dml.Clear()
	dml.SelectColumn(COLUMN_ID)
	dml.From(SCHEMA_ADMIN)
	dml.Where("", COLUMN_USER_ID, dbcore.EQUAL, user.UserId)
	dml.Where(dbcore.AND, COLUMN_PASSWORD, dbcore.EQUAL, user.Password)
	dml.OrderBy(COLUMN_ID, dbcore.ORDER_DESC)
	queryResult := dml.Execute(db.GetDb())

	user.Id = atoi(queryResult[0][COLUMN_ID])
}

func (user *Admin) Delete() {}

func (user *Admin) Update() {}

func (user *Admin) Scrape() {}

func (user *Admin) Retrieve() {}
