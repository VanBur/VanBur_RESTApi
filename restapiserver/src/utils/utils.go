/*
	Some utils to make your life easyer :)
*/
package utils

import (
	"errors"
	"fmt"
	"os/exec"
)

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
