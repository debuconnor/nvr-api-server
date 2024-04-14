package roomapi

import "github.com/debuconnor/dbcore"

func NewRoom(room Room) Website {
	return &room
}

func (room *Room) Get() {
	defer Recover()

	dml := dbcore.NewDml()
	dml.SelectAll()
	dml.From(SCHEMA_ROOM)
	dml.Where("", COLUMN_ID, dbcore.EQUAL, itoa(room.Id))
	queryResult := dml.Execute(db.GetDb())

	if len(queryResult) > 0 {
		room.Place = NewPlace(Place{
			Id: atoi(queryResult[0][COLUMN_PLACE_ID]),
		})
		room.Name = queryResult[0][COLUMN_NAME]
		room.Price = atof(queryResult[0][COLUMN_PRICE])
		room.Url = queryResult[0][COLUMN_URL]
		room.Status = Status{
			Id: atoi(queryResult[0][COLUMN_STATUS]),
		}
	} else {
		room.Id = 0
	}
}

func (room *Room) Save() {
	defer Recover()
	roomId := room.Id
	room.Get()

	if room.Id == 0 {
		dml := dbcore.NewDml()
		dml.Insert()
		dml.Into(SCHEMA_ROOM)
		dml.Value(COLUMN_ID, itoa(roomId))
		dml.Value(COLUMN_ADMIN_ID, itoa(room.Place.(*Place).Admin.(*Admin).Id))
		dml.Value(COLUMN_PLATFORM_CODE, room.Place.(*Place).Platform.(*Platform).Code)
		dml.Value(COLUMN_PLACE_ID, itoa(room.Place.(*Place).Id))
		dml.Value(COLUMN_NAME, room.Name)
		dml.Value(COLUMN_PRICE, ftoa(room.Price))
		dml.Value(COLUMN_STATUS, "0") // TODO: Set status
		dml.Value(COLUMN_DESCRIPTION, room.Description)
		dml.Value(COLUMN_URL, room.Url)
		dml.Execute(db.GetDb())
		room.Id = roomId
	} else {
		room.Update()
	}
}

func (room *Room) Delete() {}

func (room *Room) Update() {
	defer Recover()

	dml := dbcore.NewDml()
	dml.Update(SCHEMA_ROOM)
	dml.Set(COLUMN_NAME, room.Name)
	dml.Set(COLUMN_PRICE, ftoa(room.Price))
	dml.Set(COLUMN_URL, room.Url)
	dml.Set(COLUMN_DESCRIPTION, room.Description)
	dml.Where("", COLUMN_ID, dbcore.EQUAL, itoa(room.Id))
	dml.Execute(db.GetDb())
}

func (room *Room) Scrape() {}

func (room *Room) Retrieve() {}
