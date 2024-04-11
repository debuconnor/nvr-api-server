package roomapi

import (
	"fmt"

	"github.com/debuconnor/dbcore"
	"github.com/valyala/fasthttp"
)

func NewReservation(reservation Reservation) Website {
	return &reservation
}

func (reservation *Reservation) Get() {
	defer Recover()

	dml := dbcore.NewDml()
	dml.SelectAll()
	dml.From(SCHEMA_RESERVATION)
	dml.Where("", COLUMN_ID, dbcore.EQUAL, itoa(reservation.Id))
	queryResult := dml.Execute(db.GetDb())

	if len(queryResult) > 0 {
		reservation.Id = atoi(queryResult[0][COLUMN_ID])
		reservation.Admin = NewAdmin(Admin{
			Id: atoi(queryResult[0][COLUMN_ADMIN_ID]),
		})
		reservation.Customer = NewCustomer(Customer{
			Id: atoi(queryResult[0][COLUMN_CUSTOMER_ID]),
		})
		reservation.Room = NewRoom(Room{
			Id: atoi(queryResult[0][COLUMN_ROOM_ID]),
		})

		reservation.Date = queryResult[0][COLUMN_DATE]
		reservation.SpendTime = atoi(queryResult[0][COLUMN_SPEND_TIME])
		reservation.PersonCount = atoi(queryResult[0][COLUMN_PERSON_COUNT])
		reservation.Memo = queryResult[0][COLUMN_MEMO]
		reservation.CreatedAt = queryResult[0][COLUMN_CREATED_AT]
		reservation.UpdatedAt = queryResult[0][COLUMN_UPDATED_AT]

		dml.Clear()
		dml.SelectAll()
		dml.From(SCHEMA_PAYMENT)
		dml.Where("", COLUMN_ID, dbcore.EQUAL, queryResult[0][COLUMN_PAYMENT_ID])
		queryResult = dml.Execute(db.GetDb())

		reservation.Payment = Payment{
			Id:         atoi(queryResult[0][COLUMN_ID]),
			Admin:      NewAdmin(Admin{Id: atoi(queryResult[0][COLUMN_ADMIN_ID])}),
			Platform:   NewPlatform(Platform{Code: queryResult[0][COLUMN_PLATFORM_CODE]}),
			Customer:   NewCustomer(Customer{Id: atoi(queryResult[0][COLUMN_CUSTOMER_ID])}),
			Amount:     atof(queryResult[0][COLUMN_AMOUNT]),
			PaidAmount: atof(queryResult[0][COLUMN_PAID_AMOUNT]),
			PaidPoint:  atof(queryResult[0][COLUMN_PAID_POINT]),
			CreatedAt:  queryResult[0][COLUMN_CREATED_AT],
			UpdatedAt:  queryResult[0][COLUMN_UPDATED_AT],
		}
	} else {
		reservation.Id = 0
	}
}

func (reservation *Reservation) Save() {
	defer Recover()
	reservation.Room.Get()

	date := reservation.Date
	addDay, hour := convertMinuteToDayHour(itoa(reservation.SpendTime))

	shour := getHour(date)
	day := getDay(date)
	month := getMonth(date)
	year := getYear(date)
	result := false

	if isLeapYear := year%4 == 0; isLeapYear {
		MONTH_END_DAY[1] = 29
	}

	if addDay > 0 {
		for addDay > 0 {
			ehour := 0
			addDay--

			if addDay > 0 {
				ehour = 24
			} else {
				ehour = hour
			}

			sDate := itoa(year) + "-" + addDatePadding(itoa(month)) + "-" + addDatePadding(itoa(day))
			eDate := sDate
			day++

			if day > MONTH_END_DAY[month-1] {
				day -= MONTH_END_DAY[month-1]
				month++
			}

			if month > 12 {
				month -= 12
				year++
			}

			result = reservation.request(shour, ehour, sDate, eDate, STATUS_OCCUPIED)

			shour = 0

			if !result {
				break
			}
		}
	} else {
		ehour := shour + hour
		sDate := itoa(year) + "-" + addDatePadding(itoa(month)) + "-" + addDatePadding(itoa(day))
		eDate := sDate
		result = reservation.request(shour, ehour, sDate, eDate, STATUS_OCCUPIED)
	}

	if result {
		reservationId := reservation.Id
		externalPlatformCode := reservation.Platform.(*Platform).Code
		reservation.Get()

		if reservation.Id == 0 {
			dml := dbcore.NewDml()
			dml.Insert()
			dml.Into(SCHEMA_PAYMENT)
			dml.Value(COLUMN_ADMIN_ID, itoa(reservation.Admin.(*Admin).Id))
			dml.Value(COLUMN_PLATFORM_CODE, externalPlatformCode)
			dml.Value(COLUMN_CUSTOMER_ID, itoa(reservation.Customer.(*Customer).Id))
			dml.Value(COLUMN_RESERVATION_ID, itoa(reservationId))
			dml.Value(COLUMN_AMOUNT, ftoa(reservation.Payment.Amount))
			dml.Value(COLUMN_PAID_AMOUNT, ftoa(reservation.Payment.PaidAmount))
			dml.Value(COLUMN_PAID_POINT, ftoa(reservation.Payment.PaidPoint))
			dml.Value(COLUMN_STATUS, "0") // TODO: Set status
			dml.Value(COLUMN_CREATED_AT, reservation.Payment.CreatedAt)
			dml.Value(COLUMN_UPDATED_AT, reservation.Payment.UpdatedAt)
			dml.Execute(db.GetDb())

			dml.Clear()
			dml.SelectColumn(COLUMN_ID)
			dml.From(SCHEMA_PAYMENT)
			dml.Where("", COLUMN_ADMIN_ID, dbcore.EQUAL, itoa(reservation.Admin.(*Admin).Id))
			dml.Where(dbcore.AND, COLUMN_PLATFORM_CODE, dbcore.EQUAL, externalPlatformCode)
			dml.Where(dbcore.AND, COLUMN_CUSTOMER_ID, dbcore.EQUAL, itoa(reservation.Customer.(*Customer).Id))
			dml.Where(dbcore.AND, COLUMN_RESERVATION_ID, dbcore.EQUAL, itoa(reservationId))
			queryResult := dml.Execute(db.GetDb())
			reservation.Payment.Id = atoi(queryResult[0][COLUMN_ID])

			dml.Clear()
			dml.Insert()
			dml.Into(SCHEMA_RESERVATION)
			dml.Value(COLUMN_ID, itoa(reservationId))
			dml.Value(COLUMN_ADMIN_ID, itoa(reservation.Admin.(*Admin).Id))
			dml.Value(COLUMN_CUSTOMER_ID, itoa(reservation.Customer.(*Customer).Id))
			dml.Value(COLUMN_ROOM_ID, itoa(reservation.Room.(*Room).Id))
			dml.Value(COLUMN_PAYMENT_ID, itoa(reservation.Payment.Id))
			dml.Value(COLUMN_STATUS, "0") // TODO: Set status
			dml.Value(COLUMN_DATE, reservation.Date)
			dml.Value(COLUMN_SPEND_TIME, itoa(reservation.SpendTime))
			dml.Value(COLUMN_PERSON_COUNT, itoa(reservation.PersonCount))
			dml.Value(COLUMN_MEMO, reservation.Memo)
			dml.Value(COLUMN_URL, reservation.Url)
			dml.Value(COLUMN_CREATED_AT, reservation.CreatedAt)
			dml.Value(COLUMN_UPDATED_AT, reservation.UpdatedAt)
			dml.Execute(db.GetDb())
			reservation.Id = reservationId
		} else {
			reservation.Update()
		}
	}
}

func (reservation *Reservation) Delete() {
	defer Recover()

	reservation.Get()
	reservation.Room.Get()

	date := reservation.Date
	addDay, hour := convertMinuteToDayHour(itoa(reservation.SpendTime))

	shour := getHour(date)
	day := getDay(date)
	month := getMonth(date)
	year := getYear(date)
	result := false

	if isLeapYear := year%4 == 0; isLeapYear {
		MONTH_END_DAY[1] = 29
	}

	if addDay > 0 {
		for addDay > 0 {
			ehour := 0
			addDay--

			if addDay > 0 {
				ehour = 24
			} else {
				ehour = hour
			}

			sDate := itoa(year) + "-" + addDatePadding(itoa(month)) + "-" + addDatePadding(itoa(day))
			eDate := sDate
			day++

			if day > MONTH_END_DAY[month-1] {
				day -= MONTH_END_DAY[month-1]
				month++
			}

			if month > 12 {
				month -= 12
				year++
			}

			result = reservation.request(shour, ehour, sDate, eDate, STATUS_AVAILABLE)

			shour = 0
		}
	} else {
		ehour := shour + hour
		sDate := itoa(year) + "-" + addDatePadding(itoa(month)) + "-" + addDatePadding(itoa(day))
		eDate := sDate
		result = reservation.request(shour, ehour, sDate, eDate, STATUS_AVAILABLE)
	}

	if result {
		dml := dbcore.NewDml()
		dml.Delete()
		dml.From(SCHEMA_RESERVATION)
		dml.Where("", COLUMN_ID, dbcore.EQUAL, itoa(reservation.Id))
		dml.Execute(db.GetDb())

		dml.Clear()
		dml.Delete()
		dml.From(SCHEMA_PAYMENT)
		dml.Where("", COLUMN_ID, dbcore.EQUAL, itoa(reservation.Payment.Id))
		dml.Execute(db.GetDb())

		reservation = NewReservation(Reservation{}).(*Reservation)
		_ = reservation
	}
}

func (reservation *Reservation) Update() {}

func (reservation *Reservation) Parse(string) {}

func (reservation *Reservation) Scrape() {}

func (reservation *Reservation) Retrieve() {}

func (reservation *Reservation) request(shour, ehour int, sDate, eDate string, status string) bool {
	requestJson := fmt.Sprintf(`{"%s":"%s","%s":"%s","%s":"%s","%s":"%s","%s":"%s","%s":null}`,
		JSON_START_DATE, sDate,
		JSON_END_DATE, eDate,
		JSON_START_HOUR, addDatePadding(itoa(shour))+":00",
		JSON_END_HOUR, addDatePadding(itoa(ehour))+":00",
		JSON_STATUS, status,
		JSON_STOCK)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(URI_SAVE_RESERVATION_PREFIX + itoa(reservation.Room.(*Room).Place.(*Place).Id) + URI_RETRIEVE_ROOM_SUFFIX + "/" + itoa(reservation.Room.(*Room).Id) + URI_SAVE_RESERVATION_SUFFIX)
	req.Header.SetMethod(HEADER_METHOD_PATCH)
	req.Header.SetContentType("application/json; charset=UTF-8")
	req.Header.Set(HEADER_AUTHORIZATION, sessionToHeader(reservation.Platform.(*Platform).Session))
	req.Header.Set(HEADER_PLATFORM_ROLE, PLATFORM_ROLE)
	req.Header.Set(HEADER_HOST, PLATFORM_HOST)
	req.SetBodyString(requestJson)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fasthttp.Do(req, resp)
	if err != nil {
		return false
	}

	if resp.StatusCode() == fasthttp.StatusOK {
		return true
	}

	return false
}
