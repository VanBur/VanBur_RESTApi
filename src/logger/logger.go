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
	var err error
	_, err = os.Stat(config.LOG_FOLDER)
	if os.IsNotExist(err) {
		err = os.Mkdir(config.LOG_FOLDER, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	logPath := filepath.Join(config.LOG_FOLDER, config.LOG_FILE_NAME)
	var f *os.File
	f, err = os.Create(logPath)
	if err != nil {
		panic(err)
	}
	Log = log.New(f, "", log.LstdFlags|log.Lshortfile)
}
