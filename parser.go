package roomapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func parseHtml(html string) []interface{} {
	_, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		Error(err)
	}

	return nil
}

func encodeJson(data map[string]interface{}) (result string) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		Error(errors.New("ERROR_ENCODE_JSON"))
	}
	result = string(jsonStr)
	return
}

func decodeJson(jsonStr string) (result map[string]interface{}) {
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		Error(errors.New("ERROR_PARSE_JSON"))
	}
	return
}

func encodeJsonStruct(data interface{}) (result string) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		Error(errors.New("ERROR_ENCODE_JSON"))
	}
	result = string(jsonStr)
	return
}

func decodeJsonArray(jsonStr string) (result []map[string]interface{}) {
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		Error(errors.New("ERROR_PARSE_JSON"))
	}
	return
}

func atoi(str string) (result int) {
	result, err := strconv.Atoi(str)
	if err != nil {
		Error(errors.New("ERROR_PARSE_JSON"))
	}
	return
}

func atof(str string) (result float64) {
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		Error(errors.New("ERROR_PARSE_JSON"))
	}
	return
}

func itoa(num int) (result string) {
	result = strconv.Itoa(num)
	return
}

func ftoa(num float64) (result string) {
	result = strconv.FormatFloat(num, 'f', -1, 64)
	return
}

func getYear(date string) int {
	year := atoi(strings.Split(date, "-")[0])
	return year
}

func getMonth(date string) int {
	month := atoi(strings.Split(date, "-")[1])
	return month
}

func getDay(date string) int {
	dayStr := strings.Split(date, "-")[2]
	day := atoi(strings.Split(dayStr, " ")[0])

	return day
}

func getHour(date string) int {
	str := strings.Split(date, " ")[1]
	hour := atoi(strings.Split(str, ":")[0])

	return hour
}

func getNow() string {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		Error(err)
	}

	currentTime := time.Now().In(location)
	return currentTime.Format("2006-01-02 15:04:05")
}

func convertToInterfaceMap(inputMap map[string]string) map[string]interface{} {
	outputMap := make(map[string]interface{})
	for key, value := range inputMap {
		outputMap[key] = value
	}
	return outputMap
}

func convertToStringMap(inputMap map[string]interface{}) map[string]string {
	outputMap := make(map[string]string)
	for key, value := range inputMap {
		outputMap[key] = value.(string)
	}
	return outputMap
}

func convertTableColumn(table, column string) string {
	return table + "." + column
}

func parsePhoneNumber(phone string) string {
	var result []byte

	for _, v := range phone {
		if v >= '0' && v <= '9' {
			result = append(result, byte(v))
		}
	}

	return string(result)
}

func typeOf(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func convertMinuteToDayHour(minute string) (int, int) {
	minuteInt := atoi(minute)
	hourInt := minuteInt / 60
	dayInt := hourInt / 24
	hourInt = hourInt % 24

	return dayInt, hourInt
}

func addDatePadding(str string) string {
	if len(str) == 1 {
		return "0" + str
	}
	return str
}

func getJsonStringValuesByKey(jsonStr, key string) string {
	result := ""
	var bracket byte
	openBracketCount := 0
	hit := false
	key = "\"" + key + "\":"

	keyIndex := strings.Index(jsonStr, key) + len(key)

	if keyIndex > -1 {
		for i := keyIndex; i < len(jsonStr); i++ {
			if !hit {
				if jsonStr[i] == '{' || jsonStr[i] == '[' {
					openBracketCount++
					hit = true
					bracket = jsonStr[i]
				}
			} else {
				if jsonStr[i] == bracket {
					openBracketCount++
				} else if jsonStr[i] == bracket+2 {
					openBracketCount--
				}

				if openBracketCount == 0 {
					result = jsonStr[keyIndex : i+1]
					break
				}
			}
		}
	}

	return result
}
