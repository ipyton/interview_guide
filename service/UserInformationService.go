package service

import (
	"encoding/json"
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
