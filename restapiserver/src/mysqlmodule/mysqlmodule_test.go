/*
	For testing database you can use docker container with MySql server
It simple and usefull! :)

link to Docker 				: https://www.docker.com/
link to MySql Docker hub 	: https://hub.docker.com/_/mysql
*/
package mysqlmodule

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"restapiserver/src/config"
	"restapiserver/src/models"
	"restapiserver/src/utils"
)

// contentTestStruct is a structure for testing content operations
type contentTestStruct struct {
	contentData models.Content
	err         error
}

// viewContentTestStruct is a structure for testing view content operations
type viewContentTestStruct struct {
	contentData models.ViewContent
	result      string
}

var testContentData = []contentTestStruct{
	{models.Content{1, "AES 1", "mypassword", "U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"}, nil},
	{models.Content{2, "AES 2", "superpass", "U2FsdGVkX190cOearjAhFozvAQFjW53OUhLQGKfTVZnj8iOwveiaZ8rqAPNBjeDB"}, nil},
}

var testUpdatedContentData = []contentTestStruct{
	{models.Content{1, "AES 2", "", ""}, nil},
	{models.Content{2, "", "96-69", "U2FsdGVkX19wo7I1RQSYIcPzJnX8X+QXlaUSVtsXmDvzA/nxBwYQbgPpgv/CT/By"}, nil},
}

var testViewContentData = []viewContentTestStruct{
	{models.ViewContent{1, "Samsung"}, "U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"},
	{models.ViewContent{2, "LG"}, "U2FsdGVkX190cOearjAhFozvAQFjW53OUhLQGKfTVZnj8iOwveiaZ8rqAPNBjeDB"},
	{models.ViewContent{2, "iOS"}, ""},
	{models.ViewContent{999, "iOS"}, ""},
	{models.ViewContent{1, "Huawei"}, ""},
}

func TestContent(t *testing.T) {
	err := utils.StartDockerDB()
	if err != nil {
		panic(err)
	}
	defer utils.StopDockerDB()

	var dbConf = DbConfig{
		User:   utils.DOCKER_DB_USER,
		Pass:   utils.DOCKER_DB_PASS,
		Host:   utils.DOCKER_DB_HOST,
		Port:   utils.DOCKER_DB_PORT,
		DBName: utils.DOCKER_DB_NAME}
	db, err := ConnectToDataBase(dbConf)
	if err != nil {
		panic(err)
	}
	fmt.Println("DB is online")
	defer db.Close()

	err = LoadDump(db, config.DUMP_PATH)
	if err != nil {
		panic(err)
	}
	fmt.Println("Dump was loaded")

	// Get protection systems
	ps, _ := GetProtectionSystems(db)
	if len(ps) != len(models.PROTECTION_SCHEMES) {
		t.Error(
			"Expected total number of content =", len(models.PROTECTION_SCHEMES),
			"got", len(ps),
		)
	}

	// Get devices

	devices, _ := GetDevices(db)
	if len(devices) != len(models.DEVICES) {
		t.Error(
			"Expected total number of content =", len(models.DEVICES),
			"got", len(devices),
		)
	}

	// Add content tests
	for _, testPair := range testContentData {
		e := AddContent(db, testPair.contentData)
		if e != testPair.err {
			t.Error(
				"For data", testPair.contentData,
				"expected no errors, got", e,
			)
		}
	}

	// Get content by id tests
	for _, testPair := range testContentData {
		v, _ := GetContentById(db, testPair.contentData.ID)
		if v.ProtectionSystemName != testPair.contentData.ProtectionSystemName || v.Payload != testPair.contentData.Payload || v.ContentKey != testPair.contentData.ContentKey {
			t.Error(
				"For data", testPair.contentData,
				"expected data is incorrect, got", v,
			)
		}
	}

	// Get content test
	allContent, _ := GetContent(db)
	if len(allContent) != len(testContentData) {
		t.Error(
			"Expected total number of content =", len(testContentData),
			"got", len(allContent),
		)
	}

	// View content tests
	for _, testPair := range testViewContentData {
		v, e := GetEncryptedMedia(db, testPair.contentData)
		if e == nil && v.Payload != testPair.result {
			t.Error(
				"For data", testPair.contentData,
				"expected data is incorrect, got", v,
			)
		}
	}

	// Update content tests
	for _, testPair := range testUpdatedContentData {
		e := UpdateContent(db, testPair.contentData.ID, testPair.contentData)
		if e != testPair.err {
			t.Error(
				"For data", testPair.contentData,
				"expected data is incorrect, got", e,
			)
		}
	}

	for _, testPair := range testUpdatedContentData {
		v, _ := GetContentById(db, testPair.contentData.ID)
		if (testPair.contentData.ProtectionSystemName != "" && v.ProtectionSystemName != testPair.contentData.ProtectionSystemName) ||
			(testPair.contentData.Payload != "" && v.Payload != testPair.contentData.Payload) ||
			(testPair.contentData.ContentKey != "" && v.ContentKey != testPair.contentData.ContentKey) {
			t.Error(
				"For data", testPair.contentData,
				"params was not upgraded, got", v,
			)
		}
	}

	// Delete content test
	deletedId := 2
	err = DeleteContent(db, deletedId)
	if err != nil {
		t.Error(
			"Delete data was with error :", err,
			",expected", nil,
		)
	}
	v, err := GetContentById(db, deletedId)
	if err == nil {
		t.Error("Content", v, "wasn't deleted")
	}
}
