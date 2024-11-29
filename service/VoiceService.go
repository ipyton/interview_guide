package service

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"wxcloudrun-golang/db/dao"
)

var voiceDao dao.VoiceDaoImpl

func GetVoiceHandler(w http.ResponseWriter, r *http.Request) {
	question_id := r.URL.Query().Get("question_id")
	fmt.Println(question_id)
	id, err2 := strconv.ParseInt(question_id, 10, 64)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	voice, err := voiceDao.GetVoice(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		fmt.Println(err.Error())
		return
	}
	_, err = io.Copy(w, voice)
	if err != nil {
		fmt.Println(err.Error())

		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "audio/mpeg")

}

func DeleteVoiceHandler(w http.ResponseWriter, r *http.Request) {
	question_id := r.URL.Query().Get("question_id")
	err := voiceDao.DeleteVoice(question_id)
	if err != nil {
		w.Write([]byte(err.Error()))
		fmt.Println(err.Error())
		return
	}

}

func GenerateVoiceHandler(w http.ResponseWriter, r *http.Request) {
	question_id := r.URL.Query().Get("question_id")
	fmt.Println(question_id)

	voice, err := voiceDao.GenerateVoice(question_id)
	if err != nil {
		w.Write([]byte(err.Error()))
		fmt.Println(err.Error())
	}
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Write(voice)
}
