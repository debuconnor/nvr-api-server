package roomapi

import "github.com/debuconnor/dbcore"

func NewCustomer(customer Customer) User {
	return &customer
}

func (user *Customer) Get() {
	dml := dbcore.NewDml()

	if user.Phone == "" {
		dml.SelectAll()
		dml.From(SCHEMA_CUSTOMER)
		dml.Where("", COLUMN_ID, dbcore.EQUAL, itoa(user.Id))
		queryResult := dml.Execute(db.GetDb())

		user.Id = atoi(queryResult[0][COLUMN_ID])
		user.Name = queryResult[0][COLUMN_NAME]
		user.Phone = queryResult[0][COLUMN_PHONE]
		user.Email = queryResult[0][COLUMN_EMAIL]
	} else if user.Id == 0 {
		dml.SelectAll()
		dml.From(SCHEMA_CUSTOMER)
		dml.Where("", COLUMN_PHONE, dbcore.EQUAL, user.Phone)
		queryResult := dml.Execute(db.GetDb())

		user.Id = atoi(queryResult[0][COLUMN_ID])
		user.Name = queryResult[0][COLUMN_NAME]
		user.Phone = queryResult[0][COLUMN_PHONE]
		user.Email = queryResult[0][COLUMN_EMAIL]
	}
}

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
		dml.Value(COLUMN_STATUS, "0") // TODO: status code
		dml.Execute(db.GetDb())
	}
	user.Get()
}

func (user *Customer) Delete() {}

func (user *Customer) Update() {}

func (user *Customer) Scrape() {}

func (user *Customer) Retrieve() {}
