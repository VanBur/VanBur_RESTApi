package models

import "testing"

type contentTestStruct struct {
	contentData   Content
	needAllParams bool
	result        bool
}

type isAvailableDevStruct struct {
	name   string
	result *Device
}

type isAvailablePSStruct struct {
	name   string
	result *ProtectionSystem
}

type viewContentTestStruct struct {
	viewContentData ViewContent
	result          bool
}

var isAvailableDeviceTests = []isAvailableDevStruct{
	{"Android", &Device{1, "Android", 1}},
	{"Samsung", &Device{2, "Samsung", 2}},
	{"iOS", &Device{3, "iOS", 1}},
	{"Sony", nil},
	{",Yhs;__", nil},
}

var isAvailableProtectSystemsTests = []isAvailablePSStruct{
	{"AES 1", &ProtectionSystem{1, "AES 1", "AES + ECB"}},
	{"AES 2", &ProtectionSystem{2, "AES 2", "AES + CBC"}},
	{"AES + XYZ", nil},
	{"Widevine", nil},
}

var isValidContentDataTests = []contentTestStruct{
	{Content{1, "AES 1", "1234", "8765"}, true, true},
	{Content{2, "AES 1234", "", ""}, true, false},
	{Content{3, "AES 2", "", ""}, false, true},
	{Content{4, "", "xxxyyy", "43214321"}, false, true},
	{Content{5, "", "", ""}, false, false},
}

var isValidContentViewDataTests = []viewContentTestStruct{
	{ViewContent{1, "Samsung"}, true},
	{ViewContent{2, ""}, false},
	{ViewContent{3, "Huawei"}, false},
	{ViewContent{0, "LG"}, false},
}

func TestIsDeviceAvailable(t *testing.T) {
	for _, testPair := range isAvailableDeviceTests {
		v := GetDeviceByName(testPair.name)
		if v != nil && testPair.result == nil {
			t.Error(
				"For data", testPair.result,
				"expected", testPair.result,
				"got", v,
			)
		}
		if testPair.result != nil && (testPair.result.ID != v.ID || testPair.result.ProtectionSystemId != v.ProtectionSystemId || testPair.result.Name != v.Name) {
			t.Error(
				"For data", testPair.result,
				"expected", testPair.result,
				"got", v,
			)
		}
	}
}

func TestIsProtectionSchemeAvalable(t *testing.T) {
	for _, testPair := range isAvailableProtectSystemsTests {
		v := GetProtectionSchemeByName(testPair.name)
		if v != nil && testPair.result == nil {
			t.Error(
				"For data", testPair.result,
				"expected", testPair.result,
				"got", v,
			)
		}
		if testPair.result != nil && (testPair.result.ID != v.ID || testPair.result.EncryptionMode != v.EncryptionMode || testPair.result.Name != v.Name) {
			t.Error(
				"For data", testPair.result,
				"expected", testPair.result,
				"got", v,
			)
		}
	}
}

func TestIsValidContentData(t *testing.T) {
	for _, testPair := range isValidContentDataTests {
		v := IsValidContentData(testPair.contentData, testPair.needAllParams)
		if v != testPair.result {
			t.Error(
				"For data", testPair.contentData,
				"where needAllParams =", testPair.needAllParams,
				"expected", testPair.result,
				"got", v,
			)
		}
	}
}

func TestIsValidViewContentData(t *testing.T) {
	for _, testPair := range isValidContentViewDataTests {
		v := IsValidViewContentData(testPair.viewContentData)
		if v != testPair.result {
			t.Error(
				"For data", testPair.viewContentData,
				"expected", testPair.result,
				"got", v,
			)
		}
	}
}
