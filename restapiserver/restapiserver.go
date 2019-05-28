/*
	This RESTApi application example.
	Database 					: MySql
	AES-Encryption util 		: openssl
	Util for testing database 	: docker + MySql

	Created with Go ver. 1.9.2

	Please enjoy!

	Author:VanBur
*/
package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xlab/closer"
	"log"
	"os"
	"restapiserver/src/config"
	Models "restapiserver/src/models"
	MySqlModule "restapiserver/src/mysqlmodule"
	"restapiserver/src/restapi"
	"restapiserver/src/utils"
)

var db *sql.DB

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	//args parser
	appVerPtr := flag.Bool("version", false, "get current version.")
	apiPort := flag.Int("port", config.API_PORT, "port for RESTApi server")
	demoMode := flag.Bool("demo", false, "flag to launch RESTApi server in demo mode")

	dbHost := flag.String("db-host", config.DB_HOST, "MySql host address")
	dbPort := flag.Int("db-port", config.DB_PORT, "MySql port index")
	dbUser := flag.String("db-user", config.DB_USER, "MySql user ")
	dbPass := flag.String("db-pass", config.DB_PASS, "MySql password")
	dbName := flag.String("db-name", config.DB_NAME, "MySql database name")

	flag.Parse()

	if *appVerPtr {
		fmt.Println(config.VERSION)
		os.Exit(0)
	}

	fmt.Print("Checking needed utils...")
	err := utils.IsEverythingInstalled(config.REQUIRED_UTILS...)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println("Problem with utils. Error:%s", err.Error())
		utils.Log.Fatalf("Problem with utils. Error:%s", err.Error())
	}
	fmt.Println("ON")

	fmt.Print("Connecting to DataBase...")
	// Connecting to database
	var db *sql.DB
	if *demoMode {
		// In demo-mode docker will be used as mysql host
		err = utils.StartDockerDB()
		if err != nil {
			fmt.Println("ERROR")
			fmt.Printf("Problem with docker start. Error:%s", err.Error())
			utils.Log.Fatalf("Problem with docker start. Error:%s", err.Error())
		}
		defer utils.StopDockerDB()
		closer.Bind(utils.StopIfDockerStillAlive)
		var dbConf = MySqlModule.DbConfig{
			User:   utils.DOCKER_DB_USER,
			Pass:   utils.DOCKER_DB_PASS,
			Host:   utils.DOCKER_DB_HOST,
			Port:   utils.DOCKER_DB_PORT,
			DBName: utils.DOCKER_DB_NAME}
		db, err = MySqlModule.ConnectToDataBase(dbConf)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Printf("Problem with database connect. Error:%s", err.Error())
			utils.Log.Fatalf("Problem with database connect. Error:%s", err.Error())
		}
		defer db.Close()

		err = MySqlModule.LoadDump(db, "src/mysqlmodule/"+config.DUMP_PATH)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Printf("Problem with database dump. Error:%s", err.Error())
			utils.Log.Fatalf("Problem with database dump. Error:%s", err.Error())
		}
	} else {
		var dbConf = MySqlModule.DbConfig{
			User:   *dbUser,
			Pass:   *dbPass,
			Host:   *dbHost,
			Port:   *dbPort,
			DBName: *dbName}
		db, err = MySqlModule.ConnectToDataBase(dbConf)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Printf("Problem with database connect. Error:%s", err.Error())
			utils.Log.Fatalf("Problem with database connect. Error:%s", err.Error())
		}
	}
	fmt.Println("ON")

	// Update devices and protection systems cache if needed
	err = updateCacheData(db)
	if err != nil {
		fmt.Println("Trouble by updating cache.")
		utils.Log.Fatalf("'Update cache' - problem with database. Error:%s", err.Error())
	}
	restapi.RESTApi(*apiPort, db)

	closer.Hold()
}

//updateCacheData is a function for update application cached 'devices' and 'protection systems' data if needed
// return error if was some database error
func updateCacheData(db *sql.DB) error {
	protSys, err := MySqlModule.GetDataNames(db, Models.PROTECTION_SYSTEMS_TYPE)
	if err != nil {
		return err
	}
	if len(protSys) != len(Models.PROTECTION_SCHEMES) {
		Models.PROTECTION_SCHEMES = protSys
		utils.Log.Printf("Protection Systems was updated")
		fmt.Println("Protection Systems was updated")
	}

	dev, err := MySqlModule.GetDataNames(db, Models.DEVICES_TYPE)
	if err != nil {
		return err
	}
	if len(dev) != len(Models.DEVICES) {
		Models.DEVICES = dev
		utils.Log.Printf("Devices was updated")
		fmt.Println("Devices was updated")
	}
	return nil
}
