package user

import (
	"database/sql"
	"encoding/json"
	"hrm/db"
	"hrm/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//For registering a new user
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := UserModel{}
	//Parse username and password to json
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		res := middleware.Response{
			Error:   true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection :Get user password from the database
	db := db.ConnectDB()
	defer db.Close()
	stmt := `select password, username, role_id	from users 	WHERE username = $1`
	row := db.QueryRow(stmt, user.Username)
	//Create a variable pass to hold password returned from the database. It's the hashed version
	var pass string
	var roleId string
	var username string
	err := row.Scan(&pass, &username, &roleId)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		res := middleware.Response{
			Error:   true,
			Message: "Invalid Username or Password!",
		}
		json.NewEncoder(w).Encode(res)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(user.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		res := middleware.Response{
			Error:   true,
			Message: "Invalid Username or Password!",
		}
		json.NewEncoder(w).Encode(res)
	}
	//If everything is correct get use user's role_id and generate token
	token, err := middleware.GenerateJWT(username, roleId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	w.WriteHeader(http.StatusOK)
	res := middleware.Response{
		Error:   false,
		Message: token,
	}
	json.NewEncoder(w).Encode(res)
}

//For registering a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := UserModel{}
	//Parse req body to json
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//You can change bcrypt.DefaultCost to a reasonable integer like 12
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Error hashing and salting",
		}
		json.NewEncoder(w).Encode(res)
	}
	user.Password = string(hash)
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `INSERT INTO users(first_name, last_name, middle_name, username, password)
	VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(stmt, user.Firstname, user.Lastname, user.Middlename, user.Username, user.Password)
	//Checking for errors
	if err, ok := err.(*pq.Error); ok {
		//Checking for duplicate entry/unique violation
		if err.Code == "42701" || err.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			res := middleware.Response{
				Error:   true,
				Message: "User already exists",
			}
			json.NewEncoder(w).Encode(res)
		} else {
			//Checking for other errors
			w.WriteHeader(http.StatusInternalServerError)
			res := middleware.Response{
				Error:   true,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything is alright, return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error:   false,
		Message: "User added successfully",
	}
	json.NewEncoder(w).Encode(res)
}

//For changing password, either directly by the user concerned or by the admin
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Parse username and password to json
	user := UserModel{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		res := middleware.Response{
			Error:   true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//You can change bcrypt.DefaultCost to a reasonable integer like 12
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Error hashing and salting",
		}
		json.NewEncoder(w).Encode(res)
	}
	user.Password = string(hash)
	//Call db connection]
	db := db.ConnectDB()
	defer db.Close()
	stmt := `UPDATE users SET password = $2 WHERE id = $1`
	result, err := db.Exec(stmt, user.UserId, user.Password)
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check for no data found error
	if count, err := result.RowsAffected(); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		res := middleware.Response{
			Error:   true,
			Message: "Error returning rows affected by password change",
		}
		json.NewEncoder(w).Encode(res)
	} else {
		if count == 0 {
			w.WriteHeader(http.StatusNotModified)
			res := middleware.Response{
				Error:   true,
				Message: "User not found",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything went fine, return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error:   false,
		Message: "Password change successful",
	}
	json.NewEncoder(w).Encode(res)
}

//For fetching a single user
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := UserModel{}
	//Extract user id from req params
	params := mux.Vars(r)
	userId, err := strconv.Atoi(params["id"])
	if err != nil {
		res := middleware.Response{
			Error:   true,
			Message: "Unable to convert req param to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connction
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT * FROM users WHERE id = $1`
	row := db.QueryRow(stmt, userId)
	err = row.Scan(&user.UserId, &user.Firstname, &user.Middlename, &user.Lastname, &user.Username,)
	//Check if any row was returned or not
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error:   true,
			Message: "User not found",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check for other errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//if everything went well, return the user object
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := []UserModel{}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT * FROM users`
	rows, err := db.Query(stmt)
	//Check for all errors
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "P0002" || err.Code == "02000" {
			w.WriteHeader(http.StatusNotFound)
			res := middleware.Response{
				Error:   true,
				Message: "User schema not yet populated",
			}
			json.NewEncoder(w).Encode(res)
		}else {
			w.WriteHeader(http.StatusInternalServerError)
			res := middleware.Response{
				Error: true,
				Message: "Internal server error" + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	defer rows.Close()
	for rows.Next() {
		users := UserModel{}
		err := rows.Scan(
			&users.Firstname, &users.Lastname, &users.Middlename, &users.Username,
			&users.UserId,
		)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			res := middleware.Response{
				Error:   true,
				Message: "Unable to scan result set" + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
		data = append(data, users)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Extract user id from req params
	params := mux.Vars(r)
	userId, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `DELETE FROM users WHERE id = $1`
	result, err := db.Exec(stmt, userId)
	//Check if any row is affected by the delete operation
	if count, err := result.RowsAffected(); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		res := middleware.Response{
			Error:   true,
			Message: "Error counting rows affected by delete operation",
		}
		json.NewEncoder(w).Encode(res)
	} else {
		if count == 0 {
			w.WriteHeader(http.StatusNotFound)
			res := middleware.Response{
				Error:   true,
				Message: "No row affected by the delete operation",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//Check for all other errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error:   true,
			Message: "Internal server error" + err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//if everything went well, return response
	w.WriteHeader(http.StatusOK)
	res := middleware.Response{
		Error:   false,
		Message: "User removed successfully",
	}
	json.NewEncoder(w).Encode(res)
}

//For editing user: first_name, last_name, middle_name
func EditUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := UserModel{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Extract user id from req params
	params := mux.Vars(r)
	userId, err := strconv.Atoi(params["user_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//call db connect
	db := db.ConnectDB()
	defer db.Close()
	stmt := `UPDATE user SET first_name = $2, last_name = $3, middle_name = $4 
				WHERE
				 user_id = $1`
	result, err := db.Exec(stmt, uint64(userId), user.Firstname, user.Lastname, user.Middlename)
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: "Internal server error" + err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if any row was affected in the update operation
 if	count, err := result.RowsAffected(); err != nil {
		res := middleware.Response{
			Error: true,
			Message: "Error returning status of update operation",
		}
		json.NewEncoder(w).Encode(res)
	}else {
		if count == 0 {
			w.WriteHeader(http.StatusNotModified)
			res := middleware.Response{
				Error: true,
				Message: "Update operation NOT!!! successful. User probably doesnt exist.",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
//If everything was fine, return response
w.WriteHeader(http.StatusCreated)
 res := middleware.Response{
	 Error: false,
	 Message: "User object modified successfully",
 }
 json.NewEncoder(w).Encode(res)

}
