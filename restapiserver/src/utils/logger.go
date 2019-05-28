package utils

import (
	"log"
	"os"
	"restapiserver/src/config"
)

var Log *log.Logger

//initial logger method
func init() {
	if _, err := os.Stat(config.LOG_FOLDER); os.IsNotExist(err) {
		os.Mkdir(config.LOG_FOLDER, os.ModePerm)
	}
	var f, err = os.Create(config.LOG_FOLDER + "/" + config.LOG_FILE_NAME)
	if err != nil {
		panic(err)
	}
	Log = log.New(f, "", log.LstdFlags|log.Lshortfile)
}
