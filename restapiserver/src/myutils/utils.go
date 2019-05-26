/*
	Some utils to make your life easyer :)
*/
package myutils

import (
	"errors"
	"fmt"
	"os/exec"
)

//DbConfig is a structure for easyest configuring database connection
type DbConfig struct {
	User         string
	Pass         string
	Host         string
	Port         int
	DatabaseName string
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
