package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

type ClassRequest struct {
	ParentClassId int              `json:"parent_class_id,omitempty"`
	ClassId       int              `json:"class_id,omitempty"`
	Class         model.ClassModel `json:"class,omitempty"`
}

var classDao dao.ClassInterface = dao.ClassInterfaceImpl{}

func processClassInput(r *http.Request) (ClassRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ClassRequest{}, err
	}
	defer r.Body.Close()
	var requestParsed ClassRequest
	err = json.Unmarshal(body, &requestParsed)
	if err != nil {
		return ClassRequest{}, err
	}
	return requestParsed, nil

}

func UpsertClassHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only Delete method is allowed", http.StatusMethodNotAllowed)
	}
	input, err := processClassInput(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = dao.ClassInterfaceImpl{}.UpsertClass(input.Class)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteClassHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only Delete method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var parentClassId = int64(-1)
	var err error
	if r.URL.Query().Get("parent_class_id") != "" {

		parentClassId, err = strconv.ParseInt(r.URL.Query().Get("parent_class_id"), 10, 64)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = dao.ClassInterfaceImpl{}.DeleteClass(parentClassId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetClassHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var parentClassId = int64(-1)
	var err error
	if r.URL.Query().Get("parent_class_id") != "" {
		parentClassId, err = strconv.ParseInt(r.URL.Query().Get("parent_class_id"), 10, 64)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	classes, err := dao.ClassInterfaceImpl{}.GetClasses(parentClassId)
	w.Header().Set("Content-Type", "application/json")
	classesJson, err := json.Marshal(classes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(classesJson)
	//w.WriteHeader(http.StatusOK)
}
