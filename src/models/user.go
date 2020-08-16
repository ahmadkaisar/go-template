package models

import (
	"crypto/sha256"
	sql "database/sql"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strconv"

	db "../databases"
	handler "../handlers"
	middleware "../middlewares"
)

type User struct {
	Id int `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
	Email string `form:"email" json:"email"`
	Role Role `form:"role" json:"role"`
}

type UserResp struct {
	Info string `json:"info"`
	Count int `json:"count"`
	Users []User `json:"users"`
}

func (user User) Response(w http.ResponseWriter, r *http.Request, Info string, Count int, Users []User) {
	var response UserResp

	response.Info = Info
	response.Count = Count
	response.Users = Users

	handler.Header(w)
	
	json.NewEncoder(w).Encode(response)
}

func (user User) Get(w http.ResponseWriter, r *http.Request) {
	var auth Auth
	var count int
	var jwt middleware.JWT
	var users []User

	token, err := auth.Get(r)
	if err != nil {
		handler.Response(w, r, 401, "not authorized")
		return
	}

	query := db.Gorm.Table("user").Select("user.id, user.name, user.email, role.id, role.name").Joins("JOIN role ON user.role_id = role.id")
	
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
		query = query.Where("user.name LIKE ? OR user.email LIKE ? OR role.name LIKE ?", keyword, keyword, keyword)
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

	_, err = jwt.Claim(token)
	if err != nil {
		handler.Response(w, r, 500, "error on claim")
		log.Println(err)
		return
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
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Role.Id, &user.Role.Name); err != nil{
			handler.Response(w, r, 500, "scan error")
			log.Println(err.Error())
			return
		}
		users = append(users, user)
	}

	user.Response(w, r, "success", count, users)
}

func (user User) Post(w http.ResponseWriter, r *http.Request) {
	// var auth Auth
	// var jwt middleware.JWT

	// token, err := auth.Get(r)
	// if err != nil {
	// 	handler.Response(w, r, 401, "not authorized")
	// 	return
	// }

	// claims, err := jwt.Claim(token)
	// if err != nil {
	// 	handler.Response(w, r, 500, "error on claim")
	// 	log.Println(err)
	// 	return
	// } else if (int(claims["role_id"].(float64)) != 0) {
	// 	handler.Response(w, r, 403, "forbidden")
	// 	return
	// }

	stmt, err := db.SQL.Prepare("INSERT INTO user(name, email, password, role_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	name := r.FormValue("name")
	match, err := regexp.MatchString(`(^[a-zA-Z0-9[:space:]]+$)`, name)
	if !match || err != nil {
		handler.Response(w, r, 500, "invalid name")
		log.Println(err)
		return
	}

	email := r.FormValue("email")
	match, err = regexp.MatchString(`(^[a-zA-Z0-9\.\_\-]+)@([a-zA-Z0-9\.\_\-]+$)`, email)
	if !match || err != nil {
		handler.Response(w, r, 500, "invalid email")
		log.Println(err)
		return
	}

	sum := sha256.New()
	sum.Write([]byte(r.FormValue("password")))
	password_hash := hex.EncodeToString(sum.Sum(nil))
	password := password_hash

	role_id, err := strconv.Atoi(r.FormValue("role_id"))
	if err != nil {
		handler.Response(w, r, 500, "invalid role id")
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(name, email, password, role_id)
	if err != nil {
		handler.Response(w, r, 500, "stmt error")
		log.Println(err)
		return
	}

	user.Response(w, r, "success", 0, nil)
}

func (user User) GetId(w http.ResponseWriter, r *http.Request) {
	var auth Auth
	var count int
	var jwt middleware.JWT
	var rows *sql.Rows
	var users []User

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handler.Response(w, r, 500, "id not valid")
		log.Println(err)
		return
	}

	token, err := auth.Get(r)
	if err != nil {
		handler.Response(w, r, 401, "not authorized")
		return
	}

	claims, err := jwt.Claim(token)
	if err != nil {
		handler.Response(w, r, 500, "error on claim")
		log.Println(err)
		return
	} else if (int(claims["role_id"].(float64)) != 0) {
		handler.Response(w, r, 403, "forbidden")
		log.Println(err)
		return
	} else {
		rows, err = db.Gorm.Table("user").Select("user.id, user.name, user.email, role.id, role.name").Joins("JOIN role ON user.role_id = role.id").Where("user.id = ?", id).Rows()
		if err != nil {
			handler.Response(w, r, 500, "query error")
			log.Println(err)
			return
		}
	}

	defer rows.Close()

	for rows.Next() {
		if err:= rows.Scan(&user.Id, &user.Name, &user.Email, &user.Role.Id, &user.Role.Name); err != nil {
			handler.Response(w, r, 500, "scan error")
			log.Println(err)
			return
		} else {
			count = 1
			users = append(users, user)
		}
	}

	user.Response(w, r, "success", count, users)
}

func (user User) Put(w http.ResponseWriter, r *http.Request) {
	var auth Auth
	var jwt middleware.JWT

	token, err := auth.Get(r)
	if err != nil {
		handler.Response(w, r, 401, "not authorized")
		return
	}

	claims, err := jwt.Claim(token)
	if err != nil {
		handler.Response(w, r, 500, "error on claim")
		log.Println(err)
		return
	} else if (int(claims["role_id"].(float64)) != 0) {
		handler.Response(w, r, 403, "forbidden")
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handler.Response(w, r, 500, "id not valid")
		log.Println(err)
		return
	}

	stmt, err := db.SQL.Prepare("UPDATE user SET name = ?, email = ?, password = ?, role_id = ? WHERE user.id = ?")
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	name := r.FormValue("name")
	match, err := regexp.MatchString(`(^[a-zA-Z0-9[:space:]]+$)`, name)
	if !match || err != nil {
		handler.Response(w, r, 500, "invalid name")
		log.Println(err)
		return
	}

	email := r.FormValue("email")
	match, err = regexp.MatchString(`(^[a-zA-Z0-9\.\_\-]+)@([a-zA-Z0-9\.\_\-]+$)`, email)
	if !match || err != nil {
		handler.Response(w, r, 500, "invalid email")
		log.Println(err)
		return
	}

	sum := sha256.New()
	sum.Write([]byte(r.FormValue("password")))
	password_hash := hex.EncodeToString(sum.Sum(nil))
	password := password_hash
	role_id, err := strconv.Atoi(r.FormValue("role_id"))
	if err != nil {
		handler.Response(w, r, 500, "invalid role id")
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(name, email, password, role_id, id)
	if err != nil {
		handler.Response(w, r, 500, "stmt error")
		log.Println(err)
		return
	}

	user.Response(w, r, "success", 0, nil)
}

func (user User) Delete(w http.ResponseWriter, r *http.Request) {
	var auth Auth
	var jwt middleware.JWT

	token, err := auth.Get(r)
	if err != nil {
		handler.Response(w, r, 401, "not authorized")
		return
	}

	claims, err := jwt.Claim(token)
	if err != nil {
		handler.Response(w, r, 500, "error on claim")
		log.Println(err)
		return
	} else if (int(claims["role_id"].(float64)) != 0) {
		handler.Response(w, r, 403, "forbidden")
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handler.Response(w, r, 500, "id not valid")
		log.Println(err)
		return
	}

	stmt, err := db.SQL.Prepare("DELETE FROM user WHERE user.id = ?")
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

	user.Response(w, r, "success", 0, nil)
}