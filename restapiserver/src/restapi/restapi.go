package restapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	"restapiserver/src/aesmodule"
	"restapiserver/src/config"
	"restapiserver/src/models"
	"restapiserver/src/mysqlmodule"
	"restapiserver/src/utils"
)

var _db *sql.DB

func RESTApi(port int, db *sql.DB) {
	_db = db
	defer _db.Close()
	// Create some API listenners

	router := httprouter.New()
	router.GET("/content", GetContent)
	router.GET("/content/:id", GetContentById)
	router.POST("/content", AddContent)
	router.PUT("/content/:id", UpdateContent)
	router.DELETE("/content/:id", DeleteContent)
	router.POST("/content/view", CheckView)
	utils.Log.Printf("Server v%s pid=%d started with processes: %d", config.VERSION, os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

//getContent is an API GET method to getting from database content list by JSON.
/*	example of curl-request:
curl -X GET "http://localhost:5000/content"
*/
func GetContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	result, err := mysqlmodule.GetContent(_db)
	//result, err := mysqlmodule.GetContent(db)
	if err != nil {
		utils.Log.Printf("'Get content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		utils.Log.Printf("Problem with database. Error:%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	utils.Log.Printf("'Get content', ip : %s. Status : %d", r.Host, http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

//getContentById is an API GET method to getting from database content with selected ID by JSON.
// example of curl-request:
/*
curl -X GET "http://localhost:5000/content/12"
*/
func GetContentById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		utils.Log.Printf("'Get content by id', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		utils.Log.Printf("Not valid content id. Input=%s", ps.ByName("id"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := mysqlmodule.GetContentById(_db, id)
	if err != nil {
		utils.Log.Printf("'Get content by id', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusInternalServerError)
		utils.Log.Printf("Problem with database. Error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result == nil {
		utils.Log.Printf("'Get content by id', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusNotFound)
		utils.Log.Printf("Selected content wasn't found.")
		utils.Log.Printf("Input data:", id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	utils.Log.Printf("'Get content by id', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

//addContent is an API POST method to add content with selected params to database.
// example of curl-request:
/*
curl -i -H "Content-Type: application/json" -X POST
-d '{"protection_system_name":"testAes","content_key":"testKey","payload":"testPayload"}'
http://localhost:5000/content
*/
func AddContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var contentParams models.Content
	err := json.NewDecoder(r.Body).Decode(&contentParams)
	if err != nil {
		utils.Log.Printf("'Add content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		utils.Log.Printf("Not valid JSON request data. Error:%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isValid := models.IsValidContentData(contentParams, true)
	if !isValid {
		utils.Log.Printf("'Add content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		utils.Log.Printf("Not valid content data. Input data:", contentParams)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if models.IsProtectionSchemeAvalable(contentParams.ProtectionSystemName) == false {
		utils.Log.Printf("'Add content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		utils.Log.Printf("Selected Protection System isn't available. Input data:", contentParams)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = mysqlmodule.AddContent(_db, contentParams)
	if err != nil {
		utils.Log.Printf("'Add content', ip : %s. Status : %d", r.Host, http.StatusInternalServerError)
		utils.Log.Printf("Problem with database. Error:%s", err.Error())
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
func UpdateContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		utils.Log.Printf("'Update content', ip : %s. Status : %d", r.Host, http.StatusBadRequest)
		utils.Log.Printf("Not valid content id. Input=%s", ps.ByName("id"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var contentParams models.Content
	err = json.NewDecoder(r.Body).Decode(&contentParams)
	if err != nil {
		utils.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusBadRequest)
		utils.Log.Printf("Not valid JSON request data. Error:%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isValid := models.IsValidContentData(contentParams, false)
	if !isValid {
		utils.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusBadRequest)
		utils.Log.Printf("Not valid content data.Input data:", contentParams)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if contentParams.ProtectionSystemName != "" && models.IsProtectionSchemeAvalable(contentParams.ProtectionSystemName) == false {
		utils.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusBadRequest)
		utils.Log.Printf("Selected Protection System isn't available. Input data:", contentParams)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//err = mysqlmodule.UpdateContent(db, id, contentParams)
	err = mysqlmodule.UpdateContent(_db, id, contentParams)
	if err != nil {
		utils.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusInternalServerError)
		utils.Log.Printf("Problem with database. Error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.Log.Printf("'Update content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

//deleteContent is an API DELETE method to delete content with selected id from database.
// example of curl-request:
/*
curl -i -H "Content-Type: application/json" -X DELETE http://localhost:5000/content/1
*/
func DeleteContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		utils.Log.Printf("'Delete content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusBadRequest)
		utils.Log.Printf("Not valid content id.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//err = mysqlmodule.DeleteContent(db, id)
	err = mysqlmodule.DeleteContent(_db, id)
	if err != nil {
		utils.Log.Printf("'Delete content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusInternalServerError)
		utils.Log.Printf("Problem with database. Error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	utils.Log.Printf("'Delete content', id=%d, ip : %s. Status : %d", id, r.Host, http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

//checkView is an API POST method to show or not decrypted payload with selected content id and device.
// example of curl-request:
/*
curl -i -H "Content-Type: application/json" -X POST -d '{"content_id":1,"Device":"LG"}' http://localhost:5000/content/view
*/
func CheckView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var viewContentParams models.ViewContent
	err := json.NewDecoder(r.Body).Decode(&viewContentParams)
	if err != nil {
		utils.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
			viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusBadRequest)
		utils.Log.Printf("Not valid JSON request data. Error:%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isValid := models.IsValidViewContentData(viewContentParams)
	if !isValid {
		utils.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
			viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusBadRequest)
		utils.Log.Printf("Not valid view content data. Input data:", viewContentParams)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encrMedia, err := mysqlmodule.GetEncryptedMedia(_db, viewContentParams)
	if err != nil {
		utils.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
			viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusLocked)
		utils.Log.Printf("Problem with database. Error: %s", err.Error())
		w.WriteHeader(http.StatusLocked)
		return
	}

	aesOpensslType, err := models.ConvertDatabaseEncTypeToAesModule(encrMedia.EncryptionMode)
	if err != nil {
		utils.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
			viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusInternalServerError)
		utils.Log.Printf("Error:", err.Error(), ".Input data:", viewContentParams)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	decrMedia := aesmodule.Decrypter(encrMedia.ContentKey, encrMedia.Payload, aesOpensslType)
	utils.Log.Printf("'Check view content', id=%d, device=%s, ip : %s. Status : %d",
		viewContentParams.ContentID, viewContentParams.Device, r.Host, http.StatusOK)
	json.NewEncoder(w).Encode(decrMedia)
}
