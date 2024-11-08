package service

import (
	"encoding/json"
	"io"
	"net/http"
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
	input, err := processClassInput(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = dao.ClassInterfaceImpl{}.DeleteClass(input.ClassId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetClassHandler(w http.ResponseWriter, r *http.Request) {
	input, err := processClassInput(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	println(input.ParentClassId)
	classes, err := dao.ClassInterfaceImpl{}.GetClasses(input.ParentClassId)
	w.Header().Set("Content-Type", "application/json")
	classesJson, err := json.Marshal(classes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(classesJson)
	//w.WriteHeader(http.StatusOK)
}
