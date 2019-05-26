/*
	This is models to working with database data and validate this.
*/
package models

import (
	"errors"
	"restapiserver/src/aesmodule"
)

// Enum of types of data
const (
	PROTECTION_SYSTEMS_TYPE = iota
	DEVICES_TYPE
	CONTENT
)

// List of protection systems for cashe
var PROTECTION_SCHEMES = []string{
	"AES 1",
	"AES 2",
}

// List of devices for cashe
var DEVICES = []string{
	"Android",
	"Samsung",
	"iOS",
	"LG",
}

// Like a dict for connect openssl-aes-module with application
const (
	AES_ECB = "AES + ECB"
	AES_CBC = "AES + CBC"
)

// Content structure to parse data from database
type Content struct {
	ID                   int    `json:"id"`
	ProtectionSystemName string `json:"protection_system_name"`
	ContentKey           string `json:"content_key"`
	Payload              string `json:"payload"`
}

// ProtectionSystem structure to parse data from database
type ProtectionSystem struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	EncryptionMode string `json:"encryption_mode"`
}

// Device structure to parse data from database
type Device struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	ProtectionSystemName int    `json:"protection_system_name"`
}

// ViewContent structure to parse data from database
type ViewContent struct {
	ContentID int    `json:"content_id"`
	Device    string `json:"device"`
}

// EnryptedMedia structure to working with encrypted data
type EnryptedMedia struct {
	EncryptionMode string `json:"encryption_mode"`
	ContentKey     string `json:"content_key"`
	Payload        string `json:"payload"`
}

//IsValidContentData is a validator for content data
//return 'true' if content data is valid
func IsValidContentData(data Content, needAllParams bool) bool {
	if data.ContentKey == "" && data.Payload == "" && data.ProtectionSystemName == "" {
		return false
	}
	if needAllParams {
		if data.ContentKey == "" || data.Payload == "" || data.ProtectionSystemName == "" {
			return false
		}
	}
	if data.ProtectionSystemName != "" && IsProtectionSchemeAvalable(data.ProtectionSystemName) == false {
		return false
	}
	return true
}

//IsValidViewContentData is a validator for view content data
//return 'true' if view content data is valid
func IsValidViewContentData(data ViewContent) bool {
	if data.ContentID <= 0 || data.Device == "" {
		return false
	}
	if IsDeviceAvailable(data.Device) == false {
		return false
	}
	return true
}

//IsProtectionSchemeAvalable is a function that checking if selected protection system is valid
//return 'true' if selected protection system is in our database
func IsProtectionSchemeAvalable(inputType string) bool {
	for _, name := range PROTECTION_SCHEMES {
		if name == inputType {
			return true
		}
	}
	return false
}

//IsDeviceAvailable is a function that checking if selected device is valid.
//return 'true' if selected device is in our database
func IsDeviceAvailable(inputType string) bool {
	for _, name := range DEVICES {
		if name == inputType {
			return true
		}
	}
	return false
}

//ConvertDatabaseEncTypeToAesModule is a function to get understandable format for openssl-aes-module.
// return string with needed type of AES-encryption or error if openssl-aes-module cant working with this type
// of encryption
func ConvertDatabaseEncTypeToAesModule(inputType string) (string, error) {
	switch inputType {
	case AES_ECB:
		return aesmodule.TYPE_128_ECB, nil
	case AES_CBC:
		return aesmodule.TYPE_128_CBC, nil
	default:
		return "", errors.New("Selected AES encryption type is unavalable")
	}
}
