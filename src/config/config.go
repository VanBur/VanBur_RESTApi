/*
	Config file where you can find all of needed params to get this thing works.
*/
package config

//Version of application
const VERSION = "1.01"

//Log file path
const LOG_FOLDER = "logs"
const LOG_FILE_NAME = "info.log"

// MySql dump paths
const DUMP_CLEAR_DB_PATH = "test_bd_clean.sql"
const DUMP_WITH_CONTENT_PATH = "test_bd_with_content.sql"

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

// Programms required to make this super application works
var REQUIRED_UTILS = []string{
	"openssl",
	"docker",
}
