/*
	This module need for working with MySql database.
*/
package mysqlmodule

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"restapiserver/src/models"
)

// DbConfig is a structure for easyest configuring database connection
type DbConfig struct {
	User   string
	Pass   string
	Host   string
	Port   int
	DBName string
}

// params for connecting to database
const (
	_DATABASE_PING_COUNT = 4
	_DATABASE_PING_SLEEP = 5
)

// MySql commands
const (
	_GET_CONTENT_MYSQL = "SELECT content.id, protection_systems.name AS protection_system_name " +
		", content.content_key, content.payload " +
		"FROM content,protection_systems " +
		"WHERE content.protection_system_id = protection_systems.id"
	_GET_CONTENT_BY_ID_MYSQL = "SELECT content.id, protection_systems.name AS protection_system_name " +
		", content.content_key, content.payload " +
		"FROM content,protection_systems " +
		"WHERE content.protection_system_id = protection_systems.id AND content.id = ?"
	_ADD_CONTENT_MYSQL = "INSERT INTO content (protection_system_id, content_key, payload) " +
		"VALUES( (SELECT id FROM protection_systems WHERE name = ?), ?, ?)"
	_DELETE_CONTENT_MYSQL = "DELETE FROM content WHERE id = ?"

	_VIEW_CONTENT_MYSQL = "SELECT p.encryption_mode, c.content_key, c.payload " +
		"FROM protection_systems p " +
		"INNER JOIN content c ON p.id = c.protection_system_id " +
		"INNER JOIN devices d ON c.protection_system_id = d.protection_system_id " +
		"WHERE c.id = ? AND d.name = ?"

	_GET_PROTECTION_SYSTEMS = "SELECT * FROM protection_systems"
	_GET_DEVICES            = "SELECT * FROM devices"
)

var db *sql.DB

// ConnectToDataBase is a function to connect application to database and make sure that db is connected by ping.
// return pointer to database and error if something's going wrong.
//func ConnectToDataBase(cnf DbConfig) (*sql.DB, error) {
func ConnectToDataBase(cnf DbConfig) error {
	log.Println("Connecting to MySQL database")
	connSettings := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cnf.User, cnf.Pass, cnf.Host, cnf.Port, cnf.DBName)
	var err error
	db, err = sql.Open("mysql", connSettings)
	if err != nil {
		return err
	}
	for i := 0; i < _DATABASE_PING_COUNT; i++ {
		err := pingToDatabase()
		if err != nil {
			time.Sleep(_DATABASE_PING_SLEEP * time.Second)
		} else {
			break
		}
	}
	log.Println("DB connected")
	return nil
}

func DisconnectFromDataBase() {
	db.Close()
	log.Println("DB disconnected")
}

// pingToDatabase is a function to ping database for checking connect
// return error if ping was failed
func pingToDatabase() error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// GetContent is a function to get list of all content from database
// return pointer to slice of content and error if something's going wrong.
func GetContent() ([]*models.Content, error) {
	rows, err := db.Query(_GET_CONTENT_MYSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.Content, 0)
	for rows.Next() {
		content := new(models.Content)
		err := rows.Scan(&content.ID, &content.ProtectionSystemName, &content.ContentKey, &content.Payload)
		if err != nil {
			return nil, err
		}
		result = append(result, content)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// GetDevices is a function to get list of all devices data from database.
// Needs to validate cashed 'Devices' slice
// return pointer to slice of structs and error if something's going wrong.
func GetDevices() ([]models.Device, error) {
	rows, err := db.Query(_GET_DEVICES)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.Device, 0)
	for rows.Next() {
		device := new(models.Device)
		err := rows.Scan(&device.ID, &device.Name, &device.ProtectionSystemId)
		if err != nil {
			return nil, err
		}
		result = append(result, *device)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// GetProtectionSystems is a function to get list of all protection systems data from database.
// Needs to validate cashed 'Protection Systems' slice
// return pointer to slice of structs and error if something's going wrong.
func GetProtectionSystems() ([]models.ProtectionSystem, error) {
	rows, err := db.Query(_GET_PROTECTION_SYSTEMS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.ProtectionSystem, 0)
	for rows.Next() {
		ps := new(models.ProtectionSystem)
		err := rows.Scan(&ps.ID, &ps.Name, &ps.EncryptionMode)
		if err != nil {
			return nil, err
		}
		result = append(result, *ps)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// GetContentById is a function to get content data with selected id from database.
// return pointer to selected content data and error if something's going wrong.
func GetContentById(contentId int) (*models.Content, error) {
	row := db.QueryRow(_GET_CONTENT_BY_ID_MYSQL, contentId)
	content := new(models.Content)
	err := row.Scan(&content.ID, &content.ProtectionSystemName, &content.ContentKey, &content.Payload)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// AddContent is a function to add content data to database.
// return error if something's going wrong.
func AddContent(params models.Content) error {
	sql := _ADD_CONTENT_MYSQL
	_, err := prepareAndExec(sql, params.ProtectionSystemName, params.ContentKey, params.Payload)
	if err != nil {
		return err
	}
	return nil
}

// UpdateContent is a function to update content data with selected id in database.
// return error if something's going wrong.
func UpdateContent(contentId int, params models.Content) error {
	sql := generateUpdateSqlFromParams(contentId, params)
	_, err := prepareAndExec(sql)
	if err != nil {
		return err
	}
	return nil
}

// DeleteContent is a function to delete content data with selected id from database.
// return error if something's going wrong.
func DeleteContent(contentId int) error {
	sql := _DELETE_CONTENT_MYSQL
	_, err := prepareAndExec(sql, contentId)
	if err != nil {
		return err
	}
	return nil
}

// GetEncryptedMedia is a function to get enrypted data with keys from database.
// return pointer to selected enrypted data and error if something's going wrong.
func GetEncryptedMedia(params models.ViewContent) (*models.EnryptedMedia, error) {
	row := db.QueryRow(_VIEW_CONTENT_MYSQL, params.ContentID, params.Device)
	data := new(models.EnryptedMedia)
	err := row.Scan(&data.EncryptionMode, &data.ContentKey, &data.Payload)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// LoadDump is a function for loading sql-dump to database
func LoadDump(dump string) error {
	if db == nil {
		return errors.New("Database isn't connected")
	}
	file, err := os.Open(dump)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err = prepareAndExec(scanner.Text())
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// generateUpdateSqlFromParams is a function to get part of update-command.
// return string with needed parameters.
func generateUpdateSqlFromParams(contentId int, params models.Content) string {
	paramSlice := make([]string, 0)
	if params.ProtectionSystemName != "" {
		paramSlice = append(paramSlice,
			fmt.Sprintf("protection_system_id = ("+
				"SELECT id FROM protection_systems "+
				"WHERE name = '%s')", params.ProtectionSystemName))
	}
	if params.ContentKey != "" {
		paramSlice = append(paramSlice, fmt.Sprintf("content_key = '%s'", params.ContentKey))
	}
	if params.Payload != "" {
		paramSlice = append(paramSlice, fmt.Sprintf("payload = '%s'", params.Payload))
	}
	result := fmt.Sprintf("UPDATE content SET %s WHERE id = %d", strings.Join(paramSlice, ", "), contentId)
	return result
}

// prepareAndExec is a function for anti-duplication code
func prepareAndExec(sql string, execArgs ...interface{}) (sql.Result, error) {
	prepForm, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	result, err := prepForm.Exec(execArgs...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
