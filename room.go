package roomapi

import "github.com/debuconnor/dbcore"

func NewRoom(room Room) Website {
	return &room
}

func (room *Room) Get() {}

func (room *Room) Save() {
	dml := dbcore.NewDml()
	dml.Insert()
	dml.Into(SCHEMA_ROOM)
	dml.Value(COLUMN_ID, itoa(room.Id))
	dml.Value(COLUMN_ADMIN_ID, itoa(room.Place.(*Place).Admin.(*Admin).Id))
	dml.Value(COLUMN_PLATFORM_CODE, room.Place.(*Place).Platform.(*Platform).Code)
	dml.Value(COLUMN_PLACE_ID, itoa(room.Place.(*Place).Id))
	dml.Value(COLUMN_NAME, room.Name)
	dml.Value(COLUMN_PRICE, ftoa(room.Price))
	dml.Value(COLUMN_STATUS, "0") // TODO: Set status
	dml.Value(COLUMN_URL, room.Url)
	dml.Execute(db.GetDb())
}

func (room *Room) Delete() {}

func (room *Room) Update() {}

func (room *Room) Parse(string) {}

func (room *Room) Scrape() {}

func (room *Room) Retrieve() {}
