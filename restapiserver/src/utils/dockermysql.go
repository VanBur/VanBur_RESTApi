/*
	For application demo-mode and testing database you can use docker container with MySql server
It simple and usefull! :)

link to Docker 				: https://www.docker.com/
link to MySql Docker hub 	: https://hub.docker.com/_/mysql
*/
package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	// Settings for docker mysql database
	DOCKER_DB_HOST = "localhost"
	DOCKER_DB_PORT = 9999
	DOCKER_DB_USER = "root"
	DOCKER_DB_PASS = "root"
	DOCKER_DB_NAME = "testDB"

	// Docker params
	_DOCKER_CONT_NAME = "mysql-test-server"
	_DOCKER_CONT_BASE = "mysql:5.7.26"

	// Time for docker to setup mysql server (in seconds)
	_DOCKER_UPTIME = 12
)

//StartDockerDB is a method to start docker container with mysql database
func StartDockerDB() error {
	var err error
	err = IsEverythingInstalled("docker")
	if err != nil {
		return err
	}

	runArgs := []string{"run", "-d", "-p", strconv.Itoa(DOCKER_DB_PORT) + ":3306", "--name", _DOCKER_CONT_NAME,
		"--rm", "-e", "MYSQL_ROOT_PASSWORD=" + DOCKER_DB_PASS, "-e", "MYSQL_DATABASE=" + DOCKER_DB_NAME,
		"--tmpfs", "/var/lib/mysql", _DOCKER_CONT_BASE}
	err = dockerCmdExec(runArgs...)
	if err != nil {
		return err
	}
	loaderAnim(_DOCKER_UPTIME)
	return nil
}

//IsDockerRunning is a method to check - is docker alive or not
func IsDockerRunning() (bool, error) {
	out, err := exec.Command("docker", "ps", "-a", "--filter", "NAME="+_DOCKER_CONT_NAME).Output()
	if err != nil {
		return true, err
	}
	if strings.Contains(string(out), _DOCKER_CONT_NAME) {
		return true, nil
	}
	return false, nil
}

//StopIfDockerStillAlive is a method to stop docker container if it still running
func StopIfDockerStillAlive() {
	isRunning, _ := IsDockerRunning()
	if isRunning {
		fmt.Print(" Stop docker ...")
		StopDockerDB()
		fmt.Println(" Done.")
	}
}

//StopDockerDB is a method to exec stop command to docker
func StopDockerDB() {
	stopArgs := []string{"stop", _DOCKER_CONT_NAME}
	_ = dockerCmdExec(stopArgs...)
}

// dockerCmdExec is a function to execute docker commands.
// return error if something's going wrong.
func dockerCmdExec(args ...string) error {
	err := exec.Command("docker", args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// processAnim is design method to let user know that database setup is in progress
func loaderAnim(timer int) {
	for i := 0; i < timer; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print(".")
	}
}
