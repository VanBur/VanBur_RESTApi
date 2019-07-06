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
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xlab/closer"

	"restapiserver/src/config"
	dockerMYSQL "restapiserver/src/docker/mysql"
	"restapiserver/src/logger"
	"restapiserver/src/models"
	"restapiserver/src/mysqlmodule"
	"restapiserver/src/restapi"
)

// args
var (
	appVerPtr *bool
	apiPort   *int
	demoMode  *bool

	dbHost *string
	dbPort *int
	dbUser *string
	dbPass *string
	dbName *string
)

func init() {
	// args parser
	appVerPtr = flag.Bool("version", false, "get current version.")
	apiPort = flag.Int("port", config.API_PORT, "port for RESTApi server")
	demoMode = flag.Bool("demo", false, "flag to launch RESTApi server in demo mode")

	dbHost = flag.String("db-host", config.DB_HOST, "MySql host address")
	dbPort = flag.Int("db-port", config.DB_PORT, "MySql port index")
	dbUser = flag.String("db-user", config.DB_USER, "MySql user ")
	dbPass = flag.String("db-pass", config.DB_PASS, "MySql password")
	dbName = flag.String("db-name", config.DB_NAME, "MySql database name")
}

func main() {
	flag.Parse()

	// Version checker
	if *appVerPtr {
		fmt.Println(config.VERSION)
		os.Exit(0)
	}

	// Is every util in PATH
	fmt.Print("Checking needed utils...")
	err := IsEverythingInstalled(config.REQUIRED_UTILS...)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println("Problem with utils. Error:%s", err.Error())
		logger.Log.Fatalf("Problem with utils. Error:%s", err.Error())
	}
	fmt.Println("ON")

	fmt.Print("Connecting to DataBase...")
	// Connecting to database
	var db *sql.DB
	if *demoMode {
		// In demo-mode docker will be used as mysql host
		err = dockerMYSQL.StartDockerDB()
		if err != nil {
			fmt.Println("ERROR")
			fmt.Printf("Problem with docker start. Error:%s", err.Error())
			logger.Log.Fatalf("Problem with docker start. Error:%s", err.Error())
		}
		defer dockerMYSQL.StopDockerDB()
		closer.Bind(dockerMYSQL.StopIfDockerStillAlive)
		dbConfig := mysqlmodule.DbConfig{
			Host:   config.DB_HOST,
			Port:   config.DB_PORT,
			DBName: config.DB_NAME,
			User:   config.DB_USER,
			Pass:   config.DB_PASS,
		}
		db, err = mysqlmodule.ConnectToDataBase(dbConfig)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Printf("Problem with database connect. Error:%s", err.Error())
			logger.Log.Fatalf("Problem with database connect. Error:%s", err.Error())
		}
		dumpPath := filepath.Join("src", "mysqlmodule", "dumps", config.DUMP_WITH_CONTENT_PATH)
		err = mysqlmodule.LoadDump(db, dumpPath)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Printf("Problem with database dump. Error:%s", err.Error())
			logger.Log.Fatalf("Problem with database dump. Error:%s", err.Error())
		}
	} else {
		var dbConf = mysqlmodule.DbConfig{
			User:   *dbUser,
			Pass:   *dbPass,
			Host:   *dbHost,
			Port:   *dbPort,
			DBName: *dbName}
		db, err = mysqlmodule.ConnectToDataBase(dbConf)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Printf("Problem with database connect. Error:%s", err.Error())
			logger.Log.Fatalf("Problem with database connect. Error:%s", err.Error())
		}
	}
	defer db.Close()
	closer.Bind(func() {
		if db.Stats().OpenConnections != 0 {
			db.Close()
		}
	})
	fmt.Println("ON")

	// Update devices and protection systems cache if needed
	err = updateCacheData(db)
	if err != nil {
		fmt.Println("Trouble by updating cache.")
		logger.Log.Fatalf("'Update cache' - problem with database. Error:%s", err.Error())
	}
	restapi.RESTApi(*apiPort, db)

	closer.Hold()
}

//updateCacheData is a function for update application cached 'devices' and 'protection systems' data if needed
// return error if was some database error
func updateCacheData(db *sql.DB) error {
	protSys, err := mysqlmodule.GetProtectionSystems(db)
	if err != nil {
		return err
	}
	if len(protSys) != len(models.PROTECTION_SCHEMES) {
		models.PROTECTION_SCHEMES = protSys
		logger.Log.Printf("Protection Systems was updated")
		fmt.Println("Protection Systems was updated")
	}

	dev, err := mysqlmodule.GetDevices(db)
	if err != nil {
		return err
	}
	if len(dev) != len(models.DEVICES) {
		models.DEVICES = dev
		logger.Log.Printf("Devices was updated")
		fmt.Println("Devices was updated")
	}
	return nil
}

//IsEverythingInstalled is a method that checks every needed program in PATH
func IsEverythingInstalled(programs ...string) error {
	for _, util := range programs {
		_, err := exec.LookPath(util)
		if err != nil {
			return errors.New(fmt.Sprintf("Required util doesn't installed. Plaese install '%s'", util))
		}
	}
	return nil
}
