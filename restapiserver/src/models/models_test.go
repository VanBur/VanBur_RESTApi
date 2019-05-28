package models

import "testing"

type contentTestStruct struct {
	contentData   Content
	needAllParams bool
	result        bool
}

type isAvailableTestStruct struct {
	name   string
	result bool
}

type viewContentTestStruct struct {
	viewContentData ViewContent
	result          bool
}

var isAvailableDeviceTests = []isAvailableTestStruct{
	{"iOS", true},
	{"Sony", false},
	{",Yhs;__", false},
}

var isAvailableProtectSystemsTests = []isAvailableTestStruct{
	{"AES 2", true},
	{"AES + XYZ", false},
	{"Widevine", false},
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
		v := IsDeviceAvailable(testPair.name)
		if v != testPair.result {
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
		v := IsProtectionSchemeAvalable(testPair.name)
		if v != testPair.result {
			t.Error(
				"For data", testPair.name,
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
