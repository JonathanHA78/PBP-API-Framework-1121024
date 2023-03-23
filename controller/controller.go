package controller

import (
	"Explore1/model"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
)

// GetAllUsers...
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	query := "SELECT * FROM users"
	// Read from Header
	// name := r.Header.Get("Name")
	// if name!= "" {
	//	query += " WHERE name = '" + name + "'"
	//}
	name := r.URL.Query()["name"]
	age := r.URL.Query()["age"]
	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name = '" + name[0] + "'"
	}
	if age != nil {
		if name != nil {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " age='" + age[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		// send error response
		SendErrorResponse(w, "Something went wrong, please try again.")
		return
	}

	var user model.User
	var users []model.User

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &user.UserType); err != nil {
			log.Println(err)
			SendErrorResponse(w, "error result scan")
			return
		} else {
			users = append(users, user)
		}
	}

	var response model.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users

	Response(w, r, response)

}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	// Read from Request Body
	err := r.ParseForm()
	if err != nil {
		SendErrorResponse(w, "failed")
		return
	}
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	userType := r.Form.Get("user_type")

	_, errQuery := db.Exec("INSERT INTO users(name , age, address, email, password, usertype) values (?,?,?,?,?,?)", //, type ,?
		name,
		age,
		address,
		email,
		password,
		userType,
	)

	var response model.UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		fmt.Println(errQuery)
		response.Status = 400
		response.Message = "Insert Failed"
	}

	Response(w, r, response)
}

func UpdateUser(param martini.Params, w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	// Read from Request Body
	err := r.ParseForm()
	if err != nil {
		SendErrorResponse(w, "failed")
		return
	}
	// Get input from parameter
	userId := param["user_id"]

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	userType := r.Form.Get("user_type")

	sqlStatement := `
		UPDATE users 
		SET name = ?, age = ?, address =  ?, email = ?, password = ?, usertype = ?
		WHERE id = ?`

	_, errQuery := db.Exec(sqlStatement,
		name,
		age,
		address,
		email,
		password,
		userType,
		userId,
	)
	var response model.UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data.ID, _ = strconv.Atoi(userId)
		response.Data.Name = name
		response.Data.Age = age
		response.Data.Address = address
	} else {
		fmt.Println(errQuery)
		response.Status = 400
		response.Message = "Update Failed!"
	}

	Response(w, r, response)
}

func DeleteUser(param martini.Params, w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	userId := param["user_id"]

	_, errQuery := db.Exec("DELETE FROM users WHERE id=?", userId)

	var response model.UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		fmt.Println(errQuery)
		response.Status = 400
		response.Message = "Delete Failed!"
	}

	Response(w, r, response)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	// Read from Header
	platform := r.Header.Get("platform")

	// Read from Request Body
	err := r.ParseForm()
	if err != nil {
		SendErrorResponse(w, "failed")
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	query := "SELECT * FROM users WHERE email = '" + email + "' AND password = '" + password + "'"
	rows, err := db.Query(query)

	var response model.UserResponse
	if err == nil {
		var user model.User
		for rows.Next() {
			if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &user.UserType); err != nil {
				log.Println(err)
				SendErrorResponse(w, "error result scan")
				return
			}
		}
		generateToken(w, int(user.ID), user.Name, user.UserType)
		response.Status = 200
		if platform != "" {
			response.Message = "Success login from " + platform
		} else {
			response.Message = "Success login from unknown platform"
		}

	} else {
		fmt.Println(err)
		response.Status = 400
		response.Message = "Login Failed"
	}

	Response(w, r, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)

	var response model.UserResponse
	response.Status = 200
	response.Message = "Success"

	Response(w, r, response)
}
