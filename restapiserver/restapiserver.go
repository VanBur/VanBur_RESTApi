/*
	This RESTApi application example.
	Database 					: MySql
	AES-Encryption util 		: openssl
	Util for testing database 	: docker + MySql

	Created with Go ver. 1.9.2

	Please enjoy!

	Author:VanBur
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	AesModule "restapiserver/src/aesmodule"
	Config "restapiserver/src/config"
	"restapiserver/src/logger"
	Models "restapiserver/src/models"
	MySqlModule "restapiserver/src/mysqlmodule"
	MyUtils "restapiserver/src/myutils"
	"runtime"
	"strconv"
)

// Pointer to use database
var db *sql.DB

func main() {
	// Preparing part of main function
	fmt.Print("Checking needed myutils...")
	err := MyUtils.IsEverythingInstalled(Config.REQUIRED_UTILS...)
	if err != nil {
		fmt.Print("ERROR\n")
		logger.Log.Fatalf("Problem with utils. Error:%s", err.Error())
	}
	fmt.Print("ON\nConnecting to DataBase...")
	// Connecting to database
	db, err = MySqlModule.ConnectToDataBase(Config.DB_CONFIG)
	if err != nil {
		fmt.Print("ERROR\n")
		logger.Log.Fatalf("Problem with database connect. Error:%s", err.Error())
	}
	defer db.Close()
	fmt.Print("ON\n")

	// Update devices and protection systems cache if needed
	err = updateCacheData()
	if err != nil {
		fmt.Println("Trouble by updating cache.")
		logger.Log.Fatalf("'Update cache' - problem with database. Error:%s", err.Error())
		panic(err)
	}

	// Create some API listenners
	router := httprouter.New()
	router.GET("/content", getContent)
	router.GET("/content/:id", getContentById)
	router.POST("/content", addContent)
	router.PUT("/content/:id", updateContent)
	router.DELETE("/content/:id", deleteContent)
	router.POST("/content/view", checkView)
	logger.Log.Printf("Server v%s pid=%d started with processes: %d", Config.VERSION, os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))
	http.ListenAndServe(fmt.Sprintf(":%d", Config.API_PORT), router)
}

//getContent is an API GET method to getting from database content list by JSON.
/*	example of curl-request:
curl -X GET "http://localhost:5000/content"
*/
func getContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	result, err := MySqlModule.GetContent(db)
	if err != nil {
		logger.Log.Printf("'Get content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		logger.Log.Printf("Problem with database. Error:%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Log.Printf("'Get content', ip : %s. Status : %d", r.Host, http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

//getContentById is an API GET method to getting from database content with selected ID by JSON.
// example of curl-request:
/*
curl -X GET "http://localhost:5000/content/12"
*/
func getContentById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		logger.Log.Printf("'Get content by id', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		logger.Log.Printf("Not valid content id. Input=%s", ps.ByName("id"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := MySqlModule.GetContentById(db, id)
	if err != nil {
		logger.Log.Printf("'Get content by id', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusInternalServerError)
		logger.Log.Printf("Problem with database. Error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result == nil {
		logger.Log.Printf("'Get content by id', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusNotFound)
		logger.Log.Printf("Selected content wasn't found.")
		logger.Log.Printf("Input data:", id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	logger.Log.Printf("'Get content by id', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

//addContent is an API POST method to add content with selected params to database.
// example of curl-request:
/*
curl -i -H "Content-Type: application/json" -X POST
-d '{"protection_system_name":"testAes","content_key":"testKey","payload":"testPayload"}'
http://localhost:5000/content
*/
func addContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var contentParams Models.Content
	err := json.NewDecoder(r.Body).Decode(&contentParams)
	if err != nil {
		logger.Log.Printf("'Add content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		logger.Log.Printf("Not valid JSON request data. Error:%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isValid := Models.IsValidContentData(contentParams, true)
	if !isValid {
		logger.Log.Printf("'Add content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		logger.Log.Printf("Not valid content data. Input data:", contentParams)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if Models.IsProtectionSchemeAvalable(contentParams.ProtectionSystemName) == false {
		logger.Log.Printf("'Add content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		logger.Log.Printf("Selected Protection System isn't available. Input data:", contentParams)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = MySqlModule.AddContent(db, contentParams)
	if err != nil {
		logger.Log.Printf("'Add content', ip : %s. Status : %d", r.Host, http.StatusInternalServerError)
		logger.Log.Printf("Problem with database. Error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//updateContent is an API PUT method to update content with selected params in database.
// example of curl-request:
/*
curl -i -H "Content-Type: application/json" -X PUT
	-d '{"protection_system_name":"testAes","content_key":"testKey","payload":"testPayload"}'
	http://localhost:5000/content/1
*/
func updateContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		logger.Log.Printf("'Update content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		logger.Log.Printf("Not valid content id. Input=%s", ps.ByName("id"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var contentParams Models.Content
	err = json.NewDecoder(r.Body).Decode(&contentParams)
	if err != nil {
		logger.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusBadRequest)
		logger.Log.Printf("Not valid JSON request data. Error:%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isValid := Models.IsValidContentData(contentParams, false)
	if !isValid {
		logger.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusBadRequest)
		logger.Log.Printf("Not valid content data.Input data:", contentParams)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if contentParams.ProtectionSystemName != "" && Models.IsProtectionSchemeAvalable(contentParams.ProtectionSystemName) == false {
		logger.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusBadRequest)
		logger.Log.Printf("Selected Protection System isn't available. Input data:", contentParams)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = MySqlModule.UpdateContent(db, id, contentParams)
	if err != nil {
		logger.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusInternalServerError)
		logger.Log.Printf("Problem with database. Error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

//deleteContent is an API DELETE method to delete content with selected id from database.
// example of curl-request:
/*
curl -i -H "Content-Type: application/json" -X DELETE http://localhost:5000/content/1
*/
func deleteContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		logger.Log.Printf("'Delete content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusBadRequest)
		logger.Log.Printf("Not valid content id.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = MySqlModule.DeleteContent(db, id)
	if err != nil {
		logger.Log.Printf("'Delete content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusInternalServerError)
		logger.Log.Printf("Problem with database. Error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	logger.Log.Printf("'Delete content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

//checkView is an API POST method to show or not decrypted payload with selected content id and device.
// example of curl-request:
/*
curl -i -H "Content-Type: application/json" -X POST -d '{"content_id":1,"Device":"LG"}' http://localhost:5000/content/view
*/
func checkView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var viewContentParams Models.ViewContent
	err := json.NewDecoder(r.Body).Decode(&viewContentParams)
	if err != nil {
		logger.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
			viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusBadRequest)
		logger.Log.Printf("Not valid JSON request data. Error:%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isValid := Models.IsValidViewContentData(viewContentParams)
	if !isValid {
		logger.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
			viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusBadRequest)
		logger.Log.Printf("Not valid view content data. Input data:", viewContentParams)
		// not valid view content data
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encrMedia, err := MySqlModule.GetEncryptedMedia(db, viewContentParams)
	if err != nil {
		logger.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
			viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusLocked)
		logger.Log.Printf("Problem with database. Error: %s", err.Error())
		w.WriteHeader(http.StatusLocked)
		return
	}

	aesOpensslType, err := Models.ConvertDatabaseEncTypeToAesModule(encrMedia.EncryptionMode)
	if err != nil {
		logger.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
			viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusInternalServerError)
		logger.Log.Printf("Error:", err.Error(), ".Input data:", viewContentParams)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	decrMedia := AesModule.Decrypter(encrMedia.ContentKey, encrMedia.Payload, aesOpensslType)
	logger.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
		viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusOK)
	json.NewEncoder(w).Encode(decrMedia)
}

//updateCacheData is a function for update application cached 'devices' and 'protection systems' data if needed
// return error if was some database error
func updateCacheData() error {
	protSys, err := MySqlModule.GetDataNames(db, Models.PROTECTION_SYSTEMS_TYPE)
	if err != nil {
		return err
	}
	if len(protSys) != len(Models.PROTECTION_SCHEMES) {
		Models.PROTECTION_SCHEMES = protSys
		logger.Log.Printf("Protection Systems was updated")
		fmt.Println("Protection Systems was updated")
	}

	dev, err := MySqlModule.GetDataNames(db, Models.DEVICES_TYPE)
	if err != nil {
		return err
	}
	if len(dev) != len(Models.DEVICES) {
		Models.DEVICES = dev
		logger.Log.Printf("Devices was updated")
		fmt.Println("Devices was updated")
	}
	return nil
}
