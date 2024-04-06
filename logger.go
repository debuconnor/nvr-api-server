package roomapi

import (
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
)

func Log(msg ...interface{}) {
	msg = append([]interface{}{APP_NAME + ": "}, msg...)

	if isDebugMode() {
		SaveLog("", msg...)
	} else {
		log.Println(msg...)
	}
}

func Error(err error) {
	msg := append([]interface{}{"Error: "}, err)

	if isDebugMode() {
		SaveLog("error", msg)
	}

	panic(err)
}

func SaveLog(filename string, msg ...interface{}) {
	path := "./log/"

	if filename == "" {
		filename = "system"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println("Failed to create path folder:", err)
			return
		}
	}

	filePath := filepath.Join(path, filename+".log")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Failed to open error log file:", err)
		return
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(msg...)
}

func Recover() {
	err := recover()
	if err != nil {
		log.Println("Recover from ", string(debug.Stack()))
	}
}

func isDebugMode() bool {
	return os.Getenv("DEBUG_MODE") == "true"
}
