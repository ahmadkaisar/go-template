package models

import (
	sql "database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strconv"

	db "../databases"
	handler "../handlers"
)

type Role struct {
	Id int `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
}

type RoleResp struct {
	Info string `json:"info"`
	Count int `json:"count"`
	Roles []Role `json:"roles"`
}

func (role Role) Response(w http.ResponseWriter, r *http.Request, Info string, Count int, Roles []Role) {
	var response RoleResp

	response.Info = Info
	response.Count = Count
	response.Roles = Roles

	handler.Header(w)
	
	json.NewEncoder(w).Encode(response)
}

func (role Role) Get(w http.ResponseWriter, r *http.Request) {
	var count int
	var roles []Role

	query := db.Gorm.Table("role").Select("role.id, role.name")
	
	keyword, ok := r.URL.Query()["keyword"]
	if ok && len(keyword) > 0 {
		match, err := regexp.MatchString(`(^[a-zA-Z0-9[:space:]\/\-\\\|\,\.\_]+$)`, keyword[0])
		if !match || err != nil {
			handler.Response(w, r, 500, "invalid keyword")
			log.Println(err)
			return
		}
		keyword := keyword[0]
		keyword = `%` + keyword + `%`
		query = query.Where("role.name LIKE ?", keyword)
	}

	count_query := query
	page, page_ok := r.URL.Query()["page"]
	size, size_ok := r.URL.Query()["size"]
	if page_ok && size_ok {
		p, p_ok := strconv.Atoi(page[0])
		s, s_ok := strconv.Atoi(size[0])
		if p > 0 && s > 0 && p_ok == nil && s_ok == nil {
			query = query.Offset((p - 1) * s).Limit(s)
		}
	}

	count_query.Count(&count)
	rows, err := query.Rows()
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&role.Id, &role.Name); err != nil{
			handler.Response(w, r, 500, "scan error")
			log.Println(err.Error())
			return
		}
		roles = append(roles, role)
	}

	role.Response(w, r, "success", count, roles)
}

func (role Role) Post(w http.ResponseWriter, r *http.Request) {

	stmt, err := db.SQL.Prepare("INSERT INTO role(name) VALUES(?)")
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	ro := r.FormValue("name")
	match, err := regexp.MatchString(`(^[a-zA-Z0-9[:space:]]+$)`, ro)
	if !match || err != nil {
		handler.Response(w, r, 500, "invalid role")
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(ro)
	if err != nil {
		handler.Response(w, r, 500, "stmt error")
		log.Println(err)
		return
	}

	role.Response(w, r, "success", 0, nil)
}

func (role Role) GetId(w http.ResponseWriter, r *http.Request) {
	var count int
	var rows *sql.Rows
	var roles []Role

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handler.Response(w, r, 500, "id not valid")
		log.Println(err)
		return
	}

	rows, err = db.SQL.Query("SELECT role.id, role.name FROM role WHERE role.id = ?", id)
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		if err:= rows.Scan(&role.Id, &role.Name); err != nil {
			handler.Response(w, r, 500, "scan error")
			log.Println(err)
			return
		} else {
			count = 1
			roles = append(roles, role)
		}
	}

	role.Response(w, r, "success", count, roles)
}

func (role Role) Put(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handler.Response(w, r, 500, "id not valid")
		log.Println(err)
		return
	}

	stmt, err := db.SQL.Prepare("UPDATE role SET role.name = ? WHERE role.id = ?")
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	ro := r.FormValue("name")
	match, err := regexp.MatchString(`(^[a-zA-Z0-9[:space:]]+$)`, ro)
	if !match || err != nil {
		handler.Response(w, r, 500, "invalid role")
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(ro, id)
	if err != nil {
		handler.Response(w, r, 500, "stmt error")
		log.Println(err)
		return
	}

	role.Response(w, r, "success", 0, nil)
}

func (role Role) Delete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handler.Response(w, r, 500, "id not valid")
		log.Println(err)
		return
	}

	stmt, err := db.SQL.Prepare("DELETE FROM role WHERE role.id = ?")
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

	role.Response(w, r, "success", 0, nil)
}
