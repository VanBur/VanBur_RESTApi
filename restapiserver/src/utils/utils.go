/*
	Some utils to make your life easyer :)
*/
package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"restapiserver/src/aesmodule"
	"restapiserver/src/models"
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

//ConvertDatabaseEncTypeToAesModule is a function to get understandable format for openssl-aes-module.
// return string with needed type of AES-encryption or error if openssl-aes-module cant working with this type
// of encryption
func ConvertDatabaseEncTypeToAesModule(inputType string) (string, error) {
	switch inputType {
	case models.AES_ECB:
		return aesmodule.TYPE_128_ECB, nil
	case models.AES_CBC:
		return aesmodule.TYPE_128_CBC, nil
	default:
		return "", errors.New("Selected AES encryption type is unavalable")
	}
}
