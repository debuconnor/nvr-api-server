package roomapi

import (
	"sync"

	"github.com/debuconnor/dbcore"
	"github.com/valyala/fasthttp"
)

func setupHandler(ctx *fasthttp.RequestCtx) {
	defer Recover()

	dataJson := decodeJson(string(ctx.PostBody()))

	platformCode := dataJson[COLUMN_PLATFORM_CODE].(string)

	admin := NewAdmin(Admin{
		UserId:   dataJson[COLUMN_USER_ID].(string),
		Password: dataJson[COLUMN_PASSWORD].(string),
	})
	admin.Save()

	dml := dbcore.NewDml()
	dml.SelectAll()
	dml.From(SCHEMA_PLATFORM)
	dml.Where("", COLUMN_CODE, dbcore.EQUAL, platformCode)
	queryResult := dml.Execute(db.GetDb())

	platform := NewPlatform(Platform{
		Code:  queryResult[0][COLUMN_CODE],
		Admin: admin,
		Name:  queryResult[0][COLUMN_NAME],
		Url:   queryResult[0][COLUMN_URL],
	})
	platform.Retrieve()
	platform.Save()

	var places []Website
	var waitPlace sync.WaitGroup
	waitPlace.Add(len(platform.(*Platform).Places))

	for _, p := range platform.(*Platform).Places {
		go func() {
			defer waitPlace.Done()
			place := NewPlace(p)
			places = append(places, place)
			place.Retrieve()
			place.Save()
		}()
	}
	waitPlace.Wait()

	var waitPlace2 sync.WaitGroup
	waitPlace2.Add(len(platform.(*Platform).Places))
	for _, p := range places {
		go func() {
			defer waitPlace2.Done()
			var waitRoom sync.WaitGroup
			waitRoom.Add(len(p.(*Place).Rooms))

			for _, r := range p.(*Place).Rooms {
				go func() {
					defer waitRoom.Done()
					room := NewRoom(r)
					room.Save()
				}()
			}
			waitRoom.Wait()
		}()
	}
	waitPlace2.Wait()

	ctx.SetBodyString("true")
}

func addHandler(ctx *fasthttp.RequestCtx) {
	defer Recover()

	dataJson := decodeJson(string(ctx.PostBody()))

	admin := NewAdmin(Admin{
		Id: atoi(dataJson[COLUMN_ADMIN_ID].(string)),
	})

	platform := NewPlatform(Platform{
		Code:    dataJson[COLUMN_PLATFORM_FROM].(string),
		Session: getPlatformSession(admin.(*Admin).Id, dataJson[COLUMN_PLATFORM_CODE].(string)),
	})

	customer := NewCustomer(Customer{
		Name:  dataJson[COLUMN_NAME].(string),
		Phone: parsePhoneNumber(dataJson[COLUMN_PHONE].(string)),
		Email: dataJson[COLUMN_EMAIL].(string),
	})
	customer.Save()

	room := NewRoom(Room{
		Id: atoi(dataJson[COLUMN_ROOM_ID].(string)),
	})
	room.Get()

	reservation := NewReservation(Reservation{
		Id:       atoi(dataJson[COLUMN_RESERVATION_ID].(string)),
		Admin:    admin,
		Platform: platform,
		Customer: customer,
		Room:     room,
		Payment: Payment{
			Amount:     atof(dataJson[COLUMN_AMOUNT].(string)),
			PaidAmount: atof(dataJson[COLUMN_PAID_AMOUNT].(string)),
			PaidPoint:  atof(dataJson[COLUMN_PAID_POINT].(string)),
			CreatedAt:  getNow(),
			UpdatedAt:  getNow(),
		},
		Date:        dataJson[COLUMN_DATE].(string),
		SpendTime:   atoi(dataJson[COLUMN_SPEND_TIME].(string)),
		PersonCount: atoi(dataJson[COLUMN_PERSON_COUNT].(string)),
		Memo:        dataJson[COLUMN_PLATFORM_CODE].(string),
		CreatedAt:   getNow(),
		UpdatedAt:   getNow(),
	})

	reservation.Save()
	ctx.SetBodyString(itoa(reservation.(*Reservation).Id))
}

func cancelHandler(ctx *fasthttp.RequestCtx) {
	defer Recover()

	dataJson := decodeJson(string(ctx.PostBody()))

	dml := dbcore.NewDml()
	dml.SelectColumn(convertTableColumn(SCHEMA_SESSION, COLUMN_ADMIN_ID))
	dml.SelectColumn(convertTableColumn(SCHEMA_SESSION, COLUMN_PLATFORM_CODE))
	dml.SelectColumn(convertTableColumn(SCHEMA_SESSION, COLUMN_SESSION))
	dml.SelectColumn(convertTableColumn(SCHEMA_PLATFORM, COLUMN_URL))
	dml.From(SCHEMA_SESSION)
	dml.Join(dbcore.INNER_JOIN, SCHEMA_ADMIN)
	dml.On(COLUMN_ADMIN_ID, dbcore.EQUAL, COLUMN_ID)
	dml.Join(dbcore.INNER_JOIN, SCHEMA_PLATFORM)
	dml.On(COLUMN_PLATFORM_CODE, dbcore.EQUAL, COLUMN_CODE)
	dml.Where("", COLUMN_PLATFORM_CODE, dbcore.EQUAL, dataJson[COLUMN_PLATFORM_CODE].(string))
	dml.Where(dbcore.AND, COLUMN_ADMIN_ID, dbcore.EQUAL, dataJson[COLUMN_ADMIN_ID].(string))
	queryResult := dml.Execute(db.GetDb())

	admin := NewAdmin(Admin{
		Id: atoi(queryResult[0][COLUMN_ADMIN_ID]),
	})

	platform := NewPlatform(Platform{
		Code:    queryResult[0][COLUMN_CODE],
		Admin:   admin,
		Session: getPlatformSession(admin.(*Admin).Id, dataJson[COLUMN_PLATFORM_CODE].(string)),
		Url:     queryResult[0][COLUMN_URL],
	})

	reservation := NewReservation(Reservation{
		Id:       atoi(dataJson[COLUMN_RESERVATION_ID].(string)),
		Admin:    admin,
		Platform: platform,
	})

	reservation.Delete()
	ctx.SetBodyString("true")
}

func getHandler(ctx *fasthttp.RequestCtx) {
	defer Recover()

	reservationId := atoi(string(ctx.QueryArgs().Peek("reservation_id")))
	reservation := NewReservation(Reservation{
		Id: reservationId,
	})

	reservation.Get()
	ctx.SetBodyString(encodeJsonStruct(reservation))
}
