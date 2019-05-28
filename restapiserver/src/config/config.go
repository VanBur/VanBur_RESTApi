/*
	Config file where you can find all of needed params to get this thing works.
*/
package config

//Version of application
const VERSION = "1.01"

//Log file path
const LOG_FOLDER = "logs"
const LOG_FILE_NAME = "info.log"

// MySql dump path
const DUMP_PATH = "testbd.sql"

// Port for this application
var API_PORT = 5000

// Configuration of needed database
const (
	DB_HOST = "localhost"
	DB_PORT = 8889
	DB_USER = "root"
	DB_PASS = "root"
	DB_NAME = "rakuten"
)

/*
var DB_CONFIG = Utils.DbConfig{
	"root",
	"root",
	"localhost",
	8889,
	"rakuten",
}
*/

// Programms required to make this super application works
var REQUIRED_UTILS = []string{
	"openssl",
}
