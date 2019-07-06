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
	"restapiserver/src/errors"
	"restapiserver/src/logger"
	"restapiserver/src/models"
	"restapiserver/src/mysqlmodule"
)

func RESTApi(port int) {
	// Create some API listenners
	router := httprouter.New()
	router.GET("/content", GetContent)
	router.GET("/content/:id", GetContentById)
	router.POST("/content", AddContent)
	router.PUT("/content/:id", UpdateContent)
	router.DELETE("/content/:id", DeleteContent)

	router.GET("/view", ViewContent)
	logger.Log.Printf("Server v%s pid=%d started with processes: %d", config.VERSION, os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		logger.Log.Println("restapi : Port %d is already used")
		return
	}
}

// getContent is an API GET method to getting from database content list by JSON.
// example of curl-request:
// curl -X GET "http://localhost:5000/content"
func GetContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// database exec
	result, err := mysqlmodule.GetContent()
	if err != nil {
		logger.Log.Printf("restapi : Get content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadGateway, errors.DATABASE_PROBLEM)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.DATABASE_PROBLEM, http.StatusBadGateway)
		return
	}
	logger.Log.Printf("restapi : Get content : ip=%s : http status=%d", r.Host, http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		logger.Log.Println("restapi : Get content : JSON encode error :", err)
	}
}

// GetContentById is an API GET method to getting from database content with selected ID by JSON.
// example of curl-request:
// curl -X GET "http://localhost:5000/content/12"
func GetContentById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// Content_id validator
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		logger.Log.Printf("restapi : Get content by id: ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_CONTENT_ID)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_CONTENT_ID, http.StatusBadRequest)
		return
	}
	// database exec
	result, err := mysqlmodule.GetContentById(id)
	if err == sql.ErrNoRows {
		logger.Log.Printf("restapi : Get content by id : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.NO_SUCH_CONTENT_IN_DATABASE)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.NO_SUCH_CONTENT_IN_DATABASE, http.StatusNotFound)
		return
	} else if err != nil {
		logger.Log.Printf("restapi : Get content by id : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.DATABASE_PROBLEM)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.DATABASE_PROBLEM, http.StatusBadGateway)
		return
	}
	logger.Log.Printf("restapi : Get content by id : ip=%s : http status=%d", id, r.Host, http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		logger.Log.Println("restapi : Get content by id : JSON encode error :", err)
	}
}

// AddContent is an API POST method to add content with selected params to database.
// example of curl-request:
// curl -i -H "Content-Type: application/json" -X POST
// -d '{"protection_system_name":"testAes","content_key":"testKey","payload":"testPayload"}'
// http://localhost:5000/content
func AddContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// Decode input JSON
	var contentParams models.Content
	err := json.NewDecoder(r.Body).Decode(&contentParams)
	if err != nil {
		logger.Log.Printf("restapi : Add content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_JSON)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_JSON, http.StatusBadRequest)
		return
	}
	// Validate input data
	if models.IsValidContentData(contentParams, true) == false {
		logger.Log.Printf("restapi : Add content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_CONTENT_DATA)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_CONTENT_DATA, http.StatusBadRequest)
		return
	}

	// Check is inputed protection scheme is available
	if models.GetProtectionSchemeByName(contentParams.ProtectionSystemName) == nil {
		logger.Log.Printf("restapi : Add content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.NO_SUCH_PROTECTION_SYSTEM_IN_DATABASE)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.NO_SUCH_PROTECTION_SYSTEM_IN_DATABASE, http.StatusBadRequest)
		return
	}

	// Check valid input payload
	if isValidPayload(contentParams) == false {
		logger.Log.Printf("restapi : Add content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_PAYLOAD)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_PAYLOAD, http.StatusBadRequest)
		return
	}
	// database exec
	err = mysqlmodule.AddContent(contentParams)
	if err != nil {
		logger.Log.Printf("restapi : Add content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.DATABASE_PROBLEM)
		logger.Log.Println("- error info:", err.Error())
		http.Error(w, errors.DATABASE_PROBLEM, http.StatusBadGateway)
		return
	}
	logger.Log.Printf("restapi : Add content : ip=%s : http status=%d", r.Host, http.StatusOK)
	logger.Log.Println("- data:", contentParams)
	w.WriteHeader(http.StatusOK)
}

// UpdateContent is an API PUT method to update content with selected params in database.
// example of curl-request:
//
// curl -i -H "Content-Type: application/json" -X PUT
//	-d '{"protection_system_name":"testAes","content_key":"testKey","payload":"testPayload"}'
//	http://localhost:5000/content/1
func UpdateContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// Content_id validator
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		logger.Log.Printf("restapi : Update content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_CONTENT_ID)
		logger.Log.Println("- error info:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Decode input JSON
	var newContentData models.Content
	err = json.NewDecoder(r.Body).Decode(&newContentData)
	if err != nil {
		logger.Log.Printf("restapi : Update content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_JSON)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_JSON, http.StatusBadRequest)
		return
	}
	// Validate input data
	if models.IsValidContentData(newContentData, false) == false {
		logger.Log.Printf("restapi : Update content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_CONTENT_DATA)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_CONTENT_DATA, http.StatusBadRequest)
		return
	}

	// Check is inputed protection scheme is available
	if newContentData.ProtectionSystemName != "" && models.GetProtectionSchemeByName(newContentData.ProtectionSystemName) == nil {
		logger.Log.Printf("restapi : Update content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.NO_SUCH_PROTECTION_SYSTEM_IN_DATABASE)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.NO_SUCH_PROTECTION_SYSTEM_IN_DATABASE, http.StatusBadRequest)
		return
	}

	if models.IsValidContentData(newContentData, true) == false {
		oldContentData, err := mysqlmodule.GetContentById(id)
		if err == sql.ErrNoRows {
			logger.Log.Printf("restapi : Update content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.NO_SUCH_CONTENT_IN_DATABASE)
			logger.Log.Println("- error info:", err)
			http.Error(w, errors.NO_SUCH_CONTENT_IN_DATABASE, http.StatusNotFound)
			return
		} else if err != nil {
			logger.Log.Printf("restapi : Update content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.DATABASE_PROBLEM)
			logger.Log.Println("- error info:", err)
			http.Error(w, errors.DATABASE_PROBLEM, http.StatusBadGateway)
			return
		}
		err = moderateNewContentData(&newContentData, oldContentData)
		if err != nil {
			logger.Log.Printf("restapi : Add content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusExpectationFailed, errors.INVALID_PAYLOAD_IN_DATABASE)
			logger.Log.Println("- error info:", err)
			http.Error(w, errors.INVALID_PAYLOAD_IN_DATABASE, http.StatusExpectationFailed)
		}
	}

	if isValidPayload(newContentData) == false {
		logger.Log.Printf("restapi : Update content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_PAYLOAD)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_PAYLOAD, http.StatusBadRequest)
		return
	}

	// database exec
	err = mysqlmodule.UpdateContent(id, newContentData)
	if err != nil {
		logger.Log.Printf("restapi : Update content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadGateway, errors.DATABASE_PROBLEM)
		logger.Log.Println("- error info:", err.Error())
		http.Error(w, errors.DATABASE_PROBLEM, http.StatusBadGateway)
		return
	}
	logger.Log.Printf("restapi : Update content id=%d: ip=%s : http status=%d", id, r.Host, http.StatusOK)
	logger.Log.Println("- data:", newContentData)
	w.WriteHeader(http.StatusOK)
}

// DeleteContent is an API DELETE method to delete content with selected id from database.
// example of curl-request:
//
// curl -i -H "Content-Type: application/json" -X DELETE http://localhost:5000/content/1
func DeleteContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// Content_id validator
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		logger.Log.Printf("restapi : Delete content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_CONTENT_ID)
		logger.Log.Println("- error info:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// database exec
	err = mysqlmodule.DeleteContent(id)
	if err != nil {
		logger.Log.Printf("restapi : Delete content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadGateway, errors.DATABASE_PROBLEM)
		logger.Log.Println("- error info:", err.Error())
		http.Error(w, errors.DATABASE_PROBLEM, http.StatusBadGateway)
	}
	logger.Log.Printf("restapi : Delete content id=%d: ip=%s : http status=%d", id, r.Host, http.StatusOK)
	logger.Log.Println("- id:", id)
	w.WriteHeader(http.StatusOK)
}

// ViewContent is an API GET method to show or not decrypted payload with selected content id and device.
// example of curl-request:
// curl -i -H "Content-Type: application/json" -X GET -d '{"content_id":1,"Device":"LG"}' http://localhost:5000/view
func ViewContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var viewContentParams models.ViewContent
	err := json.NewDecoder(r.Body).Decode(&viewContentParams)
	if err != nil {
		logger.Log.Printf("restapi : View content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_JSON)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_JSON, http.StatusBadRequest)
		return
	}

	if models.IsValidViewContentData(viewContentParams) == false {
		logger.Log.Printf("restapi : View content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.INVALID_VIEW_CONTENT_DATA)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_VIEW_CONTENT_DATA, http.StatusBadRequest)
		return
	}

	if models.GetDeviceByName(viewContentParams.Device) == nil {
		logger.Log.Printf("restapi : View content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadRequest, errors.NO_SUCH_DEVICE_IN_DATABASE)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.NO_SUCH_DEVICE_IN_DATABASE, http.StatusBadRequest)
		return
	}

	encrMedia, err := mysqlmodule.GetEncryptedMedia(viewContentParams)
	if err == sql.ErrNoRows {
		logger.Log.Printf("restapi : View content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusUnavailableForLegalReasons, errors.CONTENT_CANT_BE_SHOWN)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.CONTENT_CANT_BE_SHOWN, http.StatusUnavailableForLegalReasons)
		return
	} else if err != nil {
		logger.Log.Printf("restapi : View content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusBadGateway, errors.DATABASE_PROBLEM)
		logger.Log.Println("- error info:", err.Error())
		http.Error(w, errors.DATABASE_PROBLEM, http.StatusBadGateway)
		return
	}

	aesOpensslType, err := aesmodule.ConvertDatabaseEncTypeToAesModule(encrMedia.EncryptionMode)
	if err != nil {
		logger.Log.Printf("restapi : View content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusUpgradeRequired, errors.NO_SUCH_PROTECTION_SYSTEM_IN_ENCRYPTION_MODULE)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.NO_SUCH_PROTECTION_SYSTEM_IN_ENCRYPTION_MODULE, http.StatusUpgradeRequired)
		return
	}
	decrMedia, err := aesmodule.Decrypter(encrMedia.ContentKey, encrMedia.Payload, aesOpensslType)
	if err != nil {
		logger.Log.Printf("restapi : View content : ip=%s : http status=%d : Error=%s", r.Host, http.StatusExpectationFailed, errors.INVALID_PAYLOAD_IN_DATABASE)
		logger.Log.Println("- error info:", err)
		http.Error(w, errors.INVALID_PAYLOAD_IN_DATABASE, http.StatusExpectationFailed)
		return
	}

	logger.Log.Printf("restapi : View content : ip=%s : http status=%d", r.Host, http.StatusOK)
	err = json.NewEncoder(w).Encode(decrMedia)
	if err != nil {
		logger.Log.Println("restapi : View content : JSON encode error :", err)
	}
}

// moderateNewContentData is a method to update old content params from database with new params
// return error if in database was badly encrypted payload
func moderateNewContentData(newContentData, oldContentData *models.Content) error {
	var encKey *string
	var prSysName *string
	if newContentData.ContentKey != "" {
		encKey = &newContentData.ContentKey
	} else {
		encKey = &oldContentData.ContentKey
	}
	if newContentData.ProtectionSystemName != "" {
		prSysName = &newContentData.ProtectionSystemName
	} else {
		prSysName = &oldContentData.ProtectionSystemName
	}

	if newContentData.Payload == "" {
		ps := models.GetProtectionSchemeByName(oldContentData.ProtectionSystemName)
		aesOpensslType, _ := aesmodule.ConvertDatabaseEncTypeToAesModule(ps.EncryptionMode)
		decodedPayload, err := aesmodule.Decrypter(oldContentData.ContentKey, oldContentData.Payload, aesOpensslType)
		if err != nil {
			return err
		}
		newPS := models.GetProtectionSchemeByName(*prSysName)
		newAesOpensslType, _ := aesmodule.ConvertDatabaseEncTypeToAesModule(newPS.EncryptionMode)
		encodedPayload, _ := aesmodule.Encrypter(*encKey, decodedPayload, newAesOpensslType)
		newContentData.Payload = encodedPayload
	}

	newContentData.ProtectionSystemName = *prSysName
	newContentData.ContentKey = *encKey
	return nil
}

// isValidPayload is a method to check is payload valid and can be decrypted
// return true if payload is valid
func isValidPayload(params models.Content) bool {
	// Check valid input payload
	ps := models.GetProtectionSchemeByName(params.ProtectionSystemName)
	if ps == nil {
		return false
	}
	aesOpensslType, err := aesmodule.ConvertDatabaseEncTypeToAesModule(ps.EncryptionMode)
	if err != nil {
		// no such AES scheme in AES module
		return false
	}
	_, err = aesmodule.Decrypter(params.ContentKey, params.Payload, aesOpensslType)
	if err != nil {
		// Bad input Payload
		return false
	}
	return true
}
