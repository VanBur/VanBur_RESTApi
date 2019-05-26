/*
	Config file where you can find all of needed params to get this thing works.
*/
package config

import Utils "rakuten/src/myutils"

//Version of application
const VERSION = "1.0"

//Log file path
const LOG_PATH = "logs/info.log"

// Port for this application
var API_PORT = 5000

// Configuration of needed database
var DB_CONFIG = Utils.DbConfig{
	"root",
	"root",
	"localhost",
	8889,
	"rakuten",
}

// Programms required to make this super application works
var REQUIRED_UTILS = []string{
	"openssl",
}
