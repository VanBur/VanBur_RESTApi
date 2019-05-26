package logger

import (
	"flag"
	"log"
	"os"
	"restapiserver/src/config"
)

var Log *log.Logger

//initial logger method
func init() {
	flag.Parse()
	var f, err = os.Create(config.LOG_PATH)
	if err != nil {
		panic(err)
	}
	Log = log.New(f, "", log.LstdFlags|log.Lshortfile)
}
