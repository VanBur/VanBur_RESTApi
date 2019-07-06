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
	"errors"
	"flag"
	"fmt"
	"log"
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
	err := IsEverythingInstalled(config.REQUIRED_UTILS...)
	if err != nil {
		log.Printf("Problem with utils. Error: %s", err.Error())
		os.Exit(0)
	}

	// Connecting to database
	if *demoMode {
		// In demo-mode docker will be used as mysql host
		err = dockerMYSQL.StartDockerDB()
		if err != nil {
			log.Printf("Problem with docker start. Error:%s", err.Error())
			logger.Log.Fatalf("Problem with docker start. Error:%s", err.Error())
		}
		defer dockerMYSQL.StopDockerDB()
		closer.Bind(dockerMYSQL.StopIfDockerStillAlive)
		dbConfig := mysqlmodule.DbConfig{
			Host:   dockerMYSQL.DOCKER_DB_HOST,
			Port:   dockerMYSQL.DOCKER_DB_PORT,
			DBName: dockerMYSQL.DOCKER_DB_NAME,
			User:   dockerMYSQL.DOCKER_DB_USER,
			Pass:   dockerMYSQL.DOCKER_DB_PASS,
		}
		err = mysqlmodule.ConnectToDataBase(dbConfig)
		if err != nil {
			log.Printf("Problem with database connect. Error: %s", err.Error())
			logger.Log.Fatalf("Problem with database connect. Error: %s", err.Error())
		}
		dumpPath := filepath.Join("src", "mysqlmodule", "dumps", config.DUMP_WITH_CONTENT_PATH)
		err = mysqlmodule.LoadDump(dumpPath)
		if err != nil {
			log.Printf("Problem with database dump. Error: %s", err.Error())
			logger.Log.Fatalf("Problem with database dump. Error: %s", err.Error())
		}
	} else {
		var dbConf = mysqlmodule.DbConfig{
			User:   *dbUser,
			Pass:   *dbPass,
			Host:   *dbHost,
			Port:   *dbPort,
			DBName: *dbName}
		err = mysqlmodule.ConnectToDataBase(dbConf)
		if err != nil {
			log.Printf("Problem with database connect. Error: %s", err.Error())
			logger.Log.Fatalf("Problem with database connect. Error: %s", err.Error())
		}
	}
	closer.Bind(mysqlmodule.DisconnectFromDataBase)
	// Update devices and protection systems cache if needed
	err = updateCacheData()
	if err != nil {
		log.Println("Trouble by updating cache.")
		logger.Log.Fatalf("'Update cache' - problem with database. Error: %s", err.Error())
	}
	restapi.RESTApi(*apiPort)

	closer.Hold()
}

//updateCacheData is a function for update application cached 'devices' and 'protection systems' data if needed
// return error if was some database error
func updateCacheData() error {
	protSys, err := mysqlmodule.GetProtectionSystems()
	if err != nil {
		return err
	}
	if len(protSys) != len(models.PROTECTION_SCHEMES) {
		models.PROTECTION_SCHEMES = protSys
		logger.Log.Printf("Protection Systems was updated")
		log.Println("Protection Systems was updated")
	}

	dev, err := mysqlmodule.GetDevices()
	if err != nil {
		return err
	}
	if len(dev) != len(models.DEVICES) {
		models.DEVICES = dev
		logger.Log.Printf("Devices was updated")
		log.Println("Devices was updated")
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
