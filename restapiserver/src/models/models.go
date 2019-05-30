/*
	This is models to working with database data and validate this.
*/
package models

// List of protection systems for cashe
var PROTECTION_SCHEMES = []ProtectionSystem{
	{1, "AES 1", "AES + ECB"},
	{2, "AES 2", "AES + CBC"},
}

// List of devices for cashe
var DEVICES = []Device{
	{1, "Android", 1},
	{2, "Samsung", 2},
	{3, "iOS", 1},
	{4, "LG", 2},
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
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	ProtectionSystemId int    `json:"protection_system_id"`
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
	// if no data
	if data.ContentKey == "" && data.Payload == "" && data.ProtectionSystemName == "" {
		return false
	}
	// if all params is required but something is gone
	if needAllParams && (data.ContentKey == "" || data.Payload == "" || data.ProtectionSystemName == "") {
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
	if GetDeviceByName(data.Device) == nil {
		return false
	}
	return true
}

//GetProtectionSchemeByName is a function for getting selected protection system data by name
// return pointer to a ProtSys if selected protection system is in our database
// and return nil if selected ProtSys is absence
func GetProtectionSchemeByName(inputType string) *ProtectionSystem {
	for _, ps := range PROTECTION_SCHEMES {
		if ps.Name == inputType {
			return &ps
		}
	}
	return nil
}

//GetDeviceByName is a function for getting selected device data by name
// return pointer to a Device if selected device is in our database
// and return nil if selected Device is absence
func GetDeviceByName(inputType string) *Device {
	for _, dev := range DEVICES {
		if dev.Name == inputType {
			return &dev
		}
	}
	return nil
}
