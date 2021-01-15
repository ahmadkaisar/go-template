package models

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	db "../databases"
	handler "../handlers"
)

type File struct {
	Id int `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
	Timestamp string `form:"timestamp" json:"timestamp"`
}

type FileResp struct {
	Info string `json:"info"`
	Count int `json:"count"`
	Files []File `json:"files"`
}

func (file File) Response(w http.ResponseWriter, r *http.Request, Info string, Count int, Files []File) {
	var response FileResp

	response.Info = Info
	response.Count = Count
	response.Files = Files

	handler.Header(w)
	
	json.NewEncoder(w).Encode(response)
}

func (file File) Get(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	filename := params["name"]
	match, err := regexp.MatchString(`(^[[:word:][:space:]\.]+$)`, filename)
	if !match || err != nil {
		handler.Response(w, r, 500, "filename not valid")
		log.Println(err)
		return
	}

	 // Check filepath
	 path := filepath.Dir("file/" + filename)
	 if path != "file" {
		handler.Response(w, r, 500, "illegal path to file")
		log.Println("illegal path to file")
		return
	 }
	 // Open file
    	f, err := os.Open("file/" + filename)
    	if err != nil {
    		handler.Response(w, r, 500, "failed to open file")
        	log.Println(err)
        	return
    	}
    	defer f.Close()

    	// Set header
    	// w.Header().Set("Content-type", "application/pdf")

    	// Stream to response
    	if _, err := io.Copy(w, f); err != nil {
    		handler.Response(w, r, 500, "failed to stream")
        	log.Println(err)
        	return
    	}
}

func (file File) Post(w http.ResponseWriter, r *http.Request) {

	// maximum 200 Mb file upload
	r.ParseMultipartForm(200 << 20)
	
	filename, err := handler.File(r)
	if err != nil {
		handler.Response(w, r, 500, err.Error())
		log.Println(err)
		return
	}

	nanotime, err := strconv.Atoi(strings.Split(filename, ".")[0])
	if err != nil {
		handler.Response(w, r, 500, "error on naming file")
		log.Println(err)
		return
	}

	timeunix := time.Unix(0, int64(nanotime))
	timestamp := timeunix.Format("02/01/2006, 15:04:05")

	stmt, err := db.SQL.Prepare("INSERT INTO file(name, timestamp) VALUES(?, ?)")
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(filename, timestamp)
	if err != nil {
		handler.Response(w, r, 500, "stmt error")
		log.Println(err)
		return
	}

	file.Response(w, r, "success", 0, nil)
}

func (file File) Put(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handler.Response(w, r, 500, "id not valid")
		log.Println(err)
		return
	}

	// maximum 200 Mb file upload
	r.ParseMultipartForm(200 << 20)

	filename, err := handler.File(r)
	if err != nil {
		handler.Response(w, r, 500, err.Error())
		log.Println(err)
		return
	}

	nanotime, err := strconv.Atoi(strings.Split(filename, ".")[0])
	if err != nil {
		handler.Response(w, r, 500, "error on naming file")
		log.Println(err)
		return
	}

	timeunix := time.Unix(0, int64(nanotime))
	timestamp := timeunix.Format("02/01/2006, 15:04:05")

	stmt, err := db.SQL.Prepare("UPDATE file SET name = ?, timestamp = ? WHERE file.id = ?")
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(filename, timestamp, id)
	if err != nil {
		handler.Response(w, r, 500, "stmt error")
		log.Println(err)
		return
	}

	file.Response(w, r, "success", 0, nil)
}

func (file File) Delete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handler.Response(w, r, 500, "id not valid")
		log.Println(err)
		return
	}

	stmt, err := db.SQL.Prepare("DELETE FROM file WHERE file.id = ?")
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		handler.Response(w, r, 500, "stmt error")
		log.Println(err)
		return
	}

	file.Response(w, r, "success", 0, nil)
}
