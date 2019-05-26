/*
	This module need for working with MySql database.
*/
package mysqlmodule

import (
	"database/sql"
	"fmt"
	Models "restapiserver/src/models"
	Utils "restapiserver/src/myutils"
	"strings"
)

// MySql commands
const (
	GET_CONTENT_MYSQL = "SELECT content.id, protection_systems.name AS protection_system_name " +
		", content.content_key, content.payload " +
		"FROM content,protection_systems " +
		"WHERE content.protection_system_id = protection_systems.id"
	GET_CONTENT_BY_ID_MYSQL = "SELECT content.id, protection_systems.name AS protection_system_name " +
		", content.content_key, content.payload " +
		"FROM content,protection_systems " +
		"WHERE content.protection_system_id = protection_systems.id AND content.id = ?"
	ADD_CONTENT_MYSQL = "INSERT INTO content (protection_system_id, content_key, payload) " +
		"VALUES( (SELECT id FROM protection_systems WHERE name = ?), ?, ?)"
	DELETE_CONTENT_MYSQL = "DELETE FROM content WHERE id = ?"

	VIEW_CONTENT_MYSQL = "SELECT p.encryption_mode, c.content_key, c.payload " +
		"FROM protection_systems p " +
		"INNER JOIN content c ON p.id = c.protection_system_id " +
		"INNER JOIN devices d ON c.protection_system_id  = d.protection_system_id " +
		"WHERE c.id = ? AND d.name = ?"

	GET_PROTECTION_SYSTEM_NAMES = "SELECT name FROM protection_systems"
	GET_DEVICE_NAMES            = "SELECT name FROM devices"
)

//ConnectToDataBase is a function to connect application to database and make sure that db is connected by ping.
// return pointer to database and error if something's going wrong.
func ConnectToDataBase(cnf Utils.DbConfig) (*sql.DB, error) {
	connSettings := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cnf.User, cnf.Pass, cnf.Host, cnf.Port, cnf.DatabaseName)
	db, err := sql.Open("mysql", connSettings)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

//GetContent is a function to get list of all content from database
// return pointer to slice of content and error if something's going wrong.
func GetContent(db *sql.DB) ([]*Models.Content, error) {
	rows, err := db.Query(GET_CONTENT_MYSQL)
	if err != nil {
		return nil, err
		//panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()

	allContent := make([]*Models.Content, 0)
	for rows.Next() {
		content := new(Models.Content)
		err := rows.Scan(&content.ID, &content.ProtectionSystemName, &content.ContentKey, &content.Payload)
		if err != nil {
			return nil, err
			//log.Fatal(err)
		}
		allContent = append(allContent, content)
	}
	if err = rows.Err(); err != nil {
		return nil, err
		//log.Fatal(err)
	}
	return allContent, nil
}

//GetDataNames is a function to get list of all names of selected data from database.
// Needs to validate cashed slices of data like 'Devices' or 'Protection systems'
// return pointer to slice of selected data names and error if something's going wrong.
func GetDataNames(db *sql.DB, dataType int) ([]string, error) {
	var sql string
	switch dataType {
	case Models.PROTECTION_SYSTEMS_TYPE:
		sql = GET_PROTECTION_SYSTEM_NAMES
	case Models.DEVICES_TYPE:
		sql = GET_DEVICE_NAMES
	}
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		result = append(result, name)
	}
	return result, nil
}

//GetContentById is a function to get content data with selected id from database.
// return pointer to selected content data and error if something's going wrong.
func GetContentById(db *sql.DB, contentId int) (*Models.Content, error) {
	row := db.QueryRow(GET_CONTENT_BY_ID_MYSQL, contentId)
	content := new(Models.Content)
	err := row.Scan(&content.ID, &content.ProtectionSystemName, &content.ContentKey, &content.Payload)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return content, nil
}

//AddContent is a function to add content data to database.
// return error if something's going wrong.
func AddContent(db *sql.DB, params Models.Content) error {
	sql := ADD_CONTENT_MYSQL
	_, err := prepareAndExec(db, sql, params.ProtectionSystemName, params.ContentKey, params.Payload)
	if err != nil {
		return err
	}
	return nil
}

//UpdateContent is a function to update content data with selected id in database.
// return error if something's going wrong.
func UpdateContent(db *sql.DB, contentId int, params Models.Content) error {
	sql := generateUpdateSqlFromParams(contentId, params)
	_, err := prepareAndExec(db, sql)
	if err != nil {
		return err
	}
	return nil
}

//DeleteContent is a function to delete content data with selected id from database.
// return error if something's going wrong.
func DeleteContent(db *sql.DB, contentId int) error {
	sql := DELETE_CONTENT_MYSQL
	_, err := prepareAndExec(db, sql, contentId)
	if err != nil {
		return err
	}
	return nil
}

//GetEncryptedMedia is a function to get enrypted data with keys from database.
// return pointer to selected enrypted data and error if something's going wrong.
func GetEncryptedMedia(db *sql.DB, params Models.ViewContent) (*Models.EnryptedMedia, error) {
	row := db.QueryRow(VIEW_CONTENT_MYSQL, params.ContentID, params.Device)
	data := new(Models.EnryptedMedia)
	err := row.Scan(&data.EncryptionMode, &data.ContentKey, &data.Payload)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return data, nil
}

//generateUpdateSqlFromParams is a function to get part of update-command.
// return string with needed parameters.
func generateUpdateSqlFromParams(contentId int, params Models.Content) string {
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

//prepareAndExec is a function for anti-duplication code
func prepareAndExec(db *sql.DB, sql string, execArgs ...interface{}) (sql.Result, error) {
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
