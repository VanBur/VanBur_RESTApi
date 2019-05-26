package mysqlmodule

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/exec"
	"rakuten/src/models"
	Utils "rakuten/src/myutils"
	"strings"
	"testing"
	"time"
)

/*
	For testing database you can use docker container with MySql server
It simple and usefull! :)

link to Docker 				: https://www.docker.com/
link to MySql Docker hub 	: https://hub.docker.com/_/mysql
*/

const (
	DOCKER_UTIL_NAME = "docker"     // Need for tests
	DELAY_FOR_MYSQL  = 10           // In secs - time to up docker with database
	DUMP_URL         = "testbd.sql" // Address of sql-dump

	// Commands for working with docker
	DOCKER_RUN = "run -d -p 9999:3306 --name=mysql-server --rm " +
		"-e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=testDB " +
		"--tmpfs /var/lib/mysql mysql:5.7.26"
	DOCKER_STOP = "stop mysql-server"
)

// Database and configs
var db *sql.DB
var test_cfg = Utils.DbConfig{
	"root",
	"root",
	"localhost",
	9999,
	"testDB",
}

// processAnim is design method to let user know that running database is in progress
func processAnim(timer int) {
	for i := 0; i < timer; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print("* ")
	}
	fmt.Println(" Done.")
}

// dockerCmdExec is a function to execute docker commands.
// return error if something's going wrong.
func dockerCmdExec(command string) error {
	args := strings.Split(command, " ")
	err := exec.Command(DOCKER_UTIL_NAME, args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// execSqlCmd is a function to execute sql commands.
// return error if something's going wrong
func execSqlCmd(cmd string) error {
	_, err := db.Exec(cmd)
	if err != nil {
		return err
	}
	return nil
}

//getEmulateTestBD is a function to read sql-dump and send commands to execute function.
// return error if something's going wrong
func getEmulateTestBD(dump string) error {
	file, err := os.Open(dump)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = execSqlCmd(scanner.Text())
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

//initTestDockerWithTestDB is a function to init test database with docker container.
// return error if something's going wrong
func initTestDockerWithTestDB() error {
	var err error
	err = Utils.IsEverythingInstalled(DOCKER_UTIL_NAME)
	if err != nil {
		return err
	}
	fmt.Println("Docker is OK")

	err = dockerCmdExec(DOCKER_RUN)
	if err != nil {
		return err
	}
	processAnim(DELAY_FOR_MYSQL)
	db, err = ConnectToDataBase(test_cfg)
	if err != nil {
		return err
	}
	fmt.Println("Docker with TestDB is up and connected")
	err = getEmulateTestBD(DUMP_URL)
	if err != nil {
		return err
	}
	return nil
}

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
	err := initTestDockerWithTestDB()
	if err != nil {
		t.Error("Problem with test Docker. Error :", err)
		return
	}
	defer dockerCmdExec(DOCKER_STOP)
	defer db.Close()

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
