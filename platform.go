package roomapi

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/debuconnor/dbcore"
	"github.com/valyala/fasthttp"
)

func NewPlatform(platform Platform) Website {
	return &platform
}

func (platform *Platform) Get() {
	defer Recover()

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
	defer Recover()
	Log("Saving platform session: ", platform.Code, ", Admin: ", platform.Admin.(*Admin).Id)

	dml := dbcore.NewDml()
	dml.Delete()
	dml.From(SCHEMA_SESSION)
	dml.Where("", COLUMN_PLATFORM_CODE, dbcore.EQUAL, platform.Code)
	dml.Where(dbcore.AND, COLUMN_ADMIN_ID, dbcore.EQUAL, itoa(platform.Admin.(*Admin).Id))
	dml.Execute(db.GetDb())

	session, err := encrypt(encodeJson(convertToInterfaceMap(platform.Session)), SECRET_SALT, SECRET_KEY)
	if err != nil {
		Error(err)
	}

	dml.Clear()
	dml.Insert()
	dml.Into(SCHEMA_SESSION)
	dml.Value(COLUMN_PLATFORM_CODE, platform.Code)
	dml.Value(COLUMN_ADMIN_ID, itoa(platform.Admin.(*Admin).Id))
	dml.Value(COLUMN_SESSION, session)
	dml.Execute(db.GetDb())
}

func (platform *Platform) Delete() {}

func (platform *Platform) Update() {}

func (platform *Platform) Parse(string) {}

func (platform *Platform) Scrape() {}

func (platform *Platform) Retrieve() {
	defer Recover()

	const isHeadless = true
	const disableExtensions = true
	const enableAutomation = false
	const timeout = 10 * time.Second
	const sleepTime = 500 * time.Millisecond

	var cookies []*network.Cookie
	var sessionCookies map[string]string

	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag(`headless`, isHeadless),
		chromedp.Flag(`disable-extensions`, disableExtensions),
		chromedp.Flag(`enable-automation`, enableAutomation),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	for retry := 0; retry < 3; retry++ {
		err := chromedp.Run(ctx,
			chromedp.Navigate(URI_RETRIEVE_SESSION),
			chromedp.ActionFunc(func(ctx context.Context) error {
				return chromedp.SetValue(SELECTOR_ID_BOX, platform.Admin.(*Admin).UserId, chromedp.NodeVisible).Do(ctx)
			}),
			chromedp.Sleep(sleepTime/2),
			chromedp.ActionFunc(func(ctx context.Context) error {
				return chromedp.SetValue(SELECTOR_PW_BOX, platform.Admin.(*Admin).Password, chromedp.NodeVisible).Do(ctx)
			}),
			chromedp.Sleep(sleepTime/2),
			chromedp.Click(SELECTOR_KEEP_LOGIN_BUTTON, chromedp.NodeVisible),
			chromedp.Click(SELECTOR_LOGIN_BUTTON, chromedp.NodeVisible),
			chromedp.Sleep(sleepTime),
			chromedp.ActionFunc(func(ctx context.Context) error {
				cookies, _ = network.GetCookies().Do(ctx)
				sessionCookies = make(map[string]string)
				for _, cookie := range cookies {
					if cookie.Name == SESSION_KEY_AUTH || cookie.Name == SESSION_KEY {
						sessionCookies[cookie.Name] = cookie.Value
					}
				}
				return nil
			}),
			chromedp.Sleep(sleepTime/2),
			chromedp.Click(SELECTOR_REGISTER_KEEPING_BUTTON, chromedp.NodeVisible),
		)

		if err == nil {
			break
		} else if err != nil && retry == 2 {
			Log(err)
		}
	}

	platform.Session = sessionCookies

	placeReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(placeReq)

	placeReq.Header.SetMethod(HEADER_METHOD_GET)
	placeReq.Header.Set(HEADER_AUTHORIZATION, sessionToHeader(platform.Session))
	placeReq.Header.Set(HEADER_PLATFORM_ID, platform.Admin.(*Admin).UserId)

	placeResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(placeResp)
	placeReq.SetRequestURI(URI_RETRIEVE_PLACE)

	err := fasthttp.Do(placeReq, placeResp)
	if err != nil {
		Error(err)
	}

	if placeResp.StatusCode() == fasthttp.StatusOK {
		placeData := decodeJsonArray(string(placeResp.Body()))
		for _, places := range placeData {
			for _, place := range places[PLATFORM_COLUMN_PLACE].([]interface{}) {
				placeMap := place.(map[string]interface{})

				place := Place{
					Id:       int(placeMap[PLATFORM_COLUMN_PLACE_ID].(float64)),
					Admin:    platform.Admin,
					Platform: platform,
					Name:     placeMap[PLATFORM_COLUMN_PLACE_NAME].(string),
					Url:      URI_RETRIEVE_ROOM_PREFIX + itoa(int(placeMap[PLATFORM_COLUMN_PLACE_ID].(float64))),
				}
				platform.Places = append(platform.Places, place)
			}
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

func sessionToHeader(session map[string]string) string {
	var header string
	for key, value := range session {
		header += key + "=" + value + "; "
	}
	return header
}
