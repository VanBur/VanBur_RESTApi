package logger

import (
	"log"
	"os"
	"path/filepath"

	"restapiserver/src/config"
)

var Log *log.Logger

//initial logger method
func init() {
	if _, err := os.Stat(config.LOG_FOLDER); os.IsNotExist(err) {
		os.Mkdir(config.LOG_FOLDER, os.ModePerm)
	}

	logPath := filepath.Join(config.LOG_FOLDER, config.LOG_FILE_NAME)
	var f, err = os.Create(logPath)
	if err != nil {
		panic(err)
	}
	Log = log.New(f, "", log.LstdFlags|log.Lshortfile)
}
