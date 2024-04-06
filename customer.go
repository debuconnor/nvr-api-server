package roomapi

import "github.com/debuconnor/dbcore"

func NewCustomer(customer Customer) User {
	return &customer
}

func (user *Customer) Get() {}

func (user *Customer) Save() {
	dml := dbcore.NewDml()
	dml.SelectAll()
	dml.From(SCHEMA_CUSTOMER)
	dml.Where("", COLUMN_PHONE, dbcore.EQUAL, user.Phone)
	queryResult := dml.Execute(db.GetDb())

	if len(queryResult) == 0 {
		dml.Clear()
		dml.Insert()
		dml.Into(SCHEMA_CUSTOMER)
		dml.Value(COLUMN_NAME, user.Name)
		dml.Value(COLUMN_PHONE, user.Phone)
		dml.Value(COLUMN_EMAIL, user.Email)
		dml.Execute(db.GetDb())
	}
}

func (user *Customer) Delete() {}

func (user *Customer) Update() {}

func (user *Customer) Scrape() {}

func (user *Customer) Retrieve() {}
