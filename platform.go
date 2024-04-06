package roomapi

import (
	"github.com/debuconnor/dbcore"
	"github.com/valyala/fasthttp"
)

func NewPlatform(platform Platform) Website {
	return &platform
}

func (platform *Platform) Get() {
	if platform.Code != "" && platform.Admin != nil {
		dml := dbcore.NewDml()
		dml.SelectColumn(convertTableColumn(SCHEMA_PLATFORM, COLUMN_NAME))
		dml.SelectColumn(convertTableColumn(SCHEMA_PLATFORM, COLUMN_URL))
		dml.SelectColumn(convertTableColumn(SCHEMA_SESSION, COLUMN_SESSION))
		dml.From(SCHEMA_SESSION)
		dml.Join(dbcore.INNER_JOIN, SCHEMA_ADMIN)
		dml.On(COLUMN_ADMIN_ID, dbcore.EQUAL, COLUMN_ID)
		dml.Join(dbcore.INNER_JOIN, SCHEMA_PLATFORM)
		dml.On(COLUMN_PLATFORM_CODE, dbcore.EQUAL, COLUMN_CODE)
		dml.Where("", COLUMN_PLATFORM_CODE, dbcore.EQUAL, platform.Code)
		dml.Where(dbcore.AND, COLUMN_ADMIN_ID, dbcore.EQUAL, itoa(platform.Admin.(*Admin).Id))
		queryResult := dml.Execute(db.GetDb())

		platform.Name = queryResult[0][COLUMN_NAME]
		platform.Url = queryResult[0][COLUMN_URL]
		platform.Session = convertToStringMap(decodeJson(queryResult[0][COLUMN_SESSION]))

		dml.Clear()
		dml.SelectColumn(convertTableColumn(SCHEMA_PLACE, COLUMN_ID))
		dml.From(SCHEMA_PLACE)
		dml.Join(dbcore.INNER_JOIN, SCHEMA_PLATFORM)
		dml.On(COLUMN_PLATFORM_CODE, dbcore.EQUAL, COLUMN_CODE)
		dml.Join(dbcore.INNER_JOIN, SCHEMA_ADMIN)
		dml.On(COLUMN_ADMIN_ID, dbcore.EQUAL, COLUMN_ID)
		dml.Where("", convertTableColumn(SCHEMA_PLATFORM, COLUMN_CODE), dbcore.EQUAL, platform.Code)
		dml.Where(dbcore.AND, convertTableColumn(SCHEMA_ADMIN, COLUMN_ID), dbcore.EQUAL, itoa(platform.Admin.(*Admin).Id))
		queryResult = dml.Execute(db.GetDb())

		for _, placeId := range queryResult {
			place := Place{
				Id: atoi(placeId[COLUMN_ID]),
			}
			platform.Places = append(platform.Places, place)
		}
	}
}

func (platform *Platform) Save() {
	dml := dbcore.NewDml()
	dml.Delete()
	dml.From(SCHEMA_SESSION)
	dml.Where("", COLUMN_PLATFORM_CODE, dbcore.EQUAL, platform.Code)
	dml.Where(dbcore.AND, COLUMN_ADMIN_ID, dbcore.EQUAL, itoa(platform.Admin.(*Admin).Id))
	dml.Execute(db.GetDb())

	dml.Clear()
	dml.Insert()
	dml.Into(SCHEMA_SESSION)
	dml.Value(COLUMN_PLATFORM_CODE, platform.Code)
	dml.Value(COLUMN_ADMIN_ID, itoa(platform.Admin.(*Admin).Id))
	dml.Value(COLUMN_SESSION, encodeJson(convertToInterfaceMap(platform.Session)))
	dml.Execute(db.GetDb())
}

func (platform *Platform) Delete() {}

func (platform *Platform) Update() {}

func (platform *Platform) Parse(string) {}

func (platform *Platform) Scrape() {}

func (platform *Platform) Retrieve() {
	sessionReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(sessionReq)
	reqBody := JSON_RETRIEVE_SESSION_PREFIX + platform.Admin.(*Admin).UserId + JSON_RETRIEVE_SESSION_MIDDLE + platform.Admin.(*Admin).Password + JSON_RETRIEVE_SESSION_SUFFIX
	sessionReq.Header.SetMethod(HEADER_METHOD_POST)
	sessionReq.Header.Set(HEADER_CONTENT_TYPE, "application/json; charset=utf-8")
	sessionReq.Header.Set(HEADER_CONTENT_LENGTH, itoa(len(reqBody)))
	sessionReq.SetBody([]byte(reqBody))

	sessionResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(sessionResp)
	sessionReq.SetRequestURI(URI_RETRIEVE_SESSION)

	err := fasthttp.Do(sessionReq, sessionResp)
	if err != nil {
		Error(err)
	}

	if sessionResp.StatusCode() == fasthttp.StatusOK {
		userData := decodeJson(string(sessionResp.Body()))
		session := userData[PLATFORM_COLUMN_USER].(map[string]interface{})[PLATFORM_COLUMN_ACCESS_TOKEN].(string)
		platform.Session = map[string]string{PLATFORM_COLUMN_ACCESS_TOKEN: session}
	}

	placeReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(placeReq)

	placeReq.Header.SetMethod(HEADER_METHOD_GET)
	placeReq.Header.Set(HEADER_AUTHORIZATION, platform.Session[PLATFORM_COLUMN_ACCESS_TOKEN])

	placeResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(placeResp)
	placeReq.SetRequestURI(URI_RETRIEVE_PLACE)

	err = fasthttp.Do(placeReq, placeResp)
	if err != nil {
		Error(err)
	}

	if placeResp.StatusCode() == fasthttp.StatusOK {
		placeData := decodeJson(string(placeResp.Body()))
		for _, place := range placeData[PLATFORM_COLUMN_PLACE].([]interface{}) {
			placeMap := place.(map[string]interface{})

			place := Place{
				Id:       int(placeMap[COLUMN_ID].(float64)),
				Admin:    platform.Admin,
				Platform: platform,
				Name:     placeMap[COLUMN_NAME].(string),
				Url:      URI_RETRIEVE_ROOM_PREFIX + itoa(int(placeMap[COLUMN_ID].(float64))) + URI_RETRIEVE_ROOM_SUFFIX,
			}
			platform.Places = append(platform.Places, place)
		}
	}
}

func getPlatformSession(adminId int, platformCode string) map[string]string {
	dml := dbcore.NewDml()
	dml.SelectColumn(COLUMN_SESSION)
	dml.From(SCHEMA_SESSION)
	dml.Where("", COLUMN_PLATFORM_CODE, dbcore.EQUAL, platformCode)
	dml.Where(dbcore.AND, COLUMN_ADMIN_ID, dbcore.EQUAL, itoa(adminId))
	queryResult := dml.Execute(db.GetDb())

	return convertToStringMap(decodeJson(queryResult[0][COLUMN_SESSION]))
}
