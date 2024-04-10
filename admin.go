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
	pw, err := decrypt(queryResult[0][COLUMN_PASSWORD], SECRET_SALT, SECRET_KEY)
	if err != nil {
		Error(err)
	} else {
		user.Password = pw
	}
}

func (user *Admin) Save() {
	exists := false

	dml := dbcore.NewDml()
	dml.SelectAll()
	dml.From(SCHEMA_ADMIN)
	dml.Where("", COLUMN_USER_ID, dbcore.EQUAL, user.UserId)
	queryResult := dml.Execute(db.GetDb())

	for _, row := range queryResult {
		ePw, _ := decrypt(row[COLUMN_PASSWORD], SECRET_SALT, SECRET_KEY)
		if row[COLUMN_USER_ID] == user.UserId && ePw == user.Password {
			exists = true
			user.Password = ePw
			break
		}
	}

	if !exists {
		pw, err := encrypt(user.Password, SECRET_SALT, SECRET_KEY)
		if err != nil {
			Error(err)
		}
		user.Password = pw
		dml.Clear()
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
		queryResult = dml.Execute(db.GetDb())
	}
	user.Id = atoi(queryResult[0][COLUMN_ID])
}

func (user *Admin) Delete() {}

func (user *Admin) Update() {}

func (user *Admin) Scrape() {}

func (user *Admin) Retrieve() {}
