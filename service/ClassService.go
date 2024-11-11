package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

var classDao dao.ClassInterface = dao.ClassInterfaceImpl{}

func processClassInput(r *http.Request) (model.ClassModel, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return model.ClassModel{}, err
	}
	defer r.Body.Close()
	var requestParsed model.ClassModel
	err = json.Unmarshal(body, &requestParsed)
	if err != nil {
		return model.ClassModel{}, err
	}
	return requestParsed, nil

}

func InsertClassHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only Post method is allowed", http.StatusMethodNotAllowed)
	}
	input, err := processClassInput(r)
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)
	}
	err = classDao.InsertClass(input)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateClassHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only Post method is allowed", http.StatusMethodNotAllowed)
	}
	input, err := processClassInput(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	err = classDao.UpdateClass(input)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteClassHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	input, err := processClassInput(r)

	if err != nil {
		fmt.Println(err.Error())
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
