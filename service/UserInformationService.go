package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"wxcloudrun-golang/db/dao"
)

type SetUserClassRequest struct {
	Openid    string `json:"openid"`
	ClassId   int    `json:"class_id"`
	ClassName string `json:"class_name"`
}

var userInterfaceDao = dao.UserInformationDaoImpl{}
var fileManagerDao = dao.FileManagerImpl{}

func SetClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var requestParsed SetUserClassRequest
	err = json.Unmarshal(body, &requestParsed)
	if err != nil {
		println(err.Error())
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	openid := r.Header.Get("openid")
	err = userInterfaceDao.UpdateUserClass(openid, requestParsed.ClassId, requestParsed.ClassName)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	openid := r.Header.Get("openid")
	info, err := userInterfaceDao.GetUserInfo(openid)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := JsonResult{}
	res.Code = 1
	res.Data = info
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)

}

func GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	openid := r.Header.Get("openid")
	file, err := fileManagerDao.GetFile("/"+openid[0:2]+"/"+openid, "avatar")
	fmt.Println(file)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// Parse the form data, which allows us to access the uploaded file
	openid := r.Header.Get("openid")
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}

	// Generate a unique filename (you can use UUID or timestamp)
	fileManagerDao.UploadFileByMultipartFile(openid, "avatar", file)
	// Upload the file to MinIO
	if err != nil {
		http.Error(w, "Failed to upload image to MinIO", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.Write([]byte("Image uploaded successfully"))

}

type Rename struct {
	Name string `json:"username"`
}

func SetUserName(w http.ResponseWriter, r *http.Request) {
	openid := r.Header.Get("openid")
	rename := Rename{}
	all, _ := io.ReadAll(r.Body)
	json.Unmarshal(all, &rename)
	userInterfaceDao.UpdateUserName(openid, rename.Name)

}
