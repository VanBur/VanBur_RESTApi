/*
	For application demo-mode and testing database you can use docker container with MySql server
It simple and usefull! :)

link to Docker 				: https://www.docker.com/
link to MySql Docker hub 	: https://hub.docker.com/_/mysql
*/
package mysql

import (
	"log"
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
	DOCKER_CONT_NAME = "mysql-test-server"
	DOCKER_CONT_BASE = "mysql:5.7.26"

	// Time for docker to setup mysql server (in seconds)
	DOCKER_UPTIME = 10
)

// StartDockerDB is a method to start docker container with mysql database
func StartDockerDB() error {
	log.Println("Start Docker container")
	var err error
	runArgs := []string{
		"run",
		"-d",
		"-p", strconv.Itoa(DOCKER_DB_PORT) + ":3306",
		"--name", DOCKER_CONT_NAME,
		"--rm",
		"-e", "MYSQL_ROOT_PASSWORD=" + DOCKER_DB_PASS,
		"-e", "MYSQL_DATABASE=" + DOCKER_DB_NAME,
		"--tmpfs", "/var/lib/mysql", DOCKER_CONT_BASE}
	err = exec.Command("docker", runArgs...).Run()
	if err != nil {
		return err
	}
	log.Printf("Docker container is running. MySQL database need %d sec. Please wait.", DOCKER_UPTIME)
	time.Sleep(DOCKER_UPTIME * time.Second)
	log.Println("Docker container is ready.")
	return nil
}

// IsDockerRunning is a method to check - is docker alive or not
func IsDockerRunning() (bool, error) {
	out, err := exec.Command("docker", "ps", "-a",
		"--filter", "NAME="+DOCKER_CONT_NAME).Output()
	if err != nil {
		return true, err
	}
	if strings.Contains(string(out), DOCKER_CONT_NAME) {
		return true, nil
	}
	return false, nil
}

// StopIfDockerStillAlive is a method to stop docker container if it still running
func StopIfDockerStillAlive() {
	isRunning, _ := IsDockerRunning()
	if isRunning {
		log.Println("Docker container still running, stopping docker.")
		StopDockerDB()
	}
}

//StopDockerDB is a method to exec stop command to docker
func StopDockerDB() {
	stopArgs := []string{"stop", DOCKER_CONT_NAME}
	err := exec.Command("docker", stopArgs...).Run()
	if err != nil {
		log.Println(err)
	}
	log.Println("Docker container stopped.")
}
