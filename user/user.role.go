package user

import (
	"database/sql"
	"encoding/json"
	"hrm/db"
	"hrm/middleware"
	"hrm/role"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

//It's an update operation that updates role_id on users table
func AssignRoleToUser(w http.ResponseWriter, r *http.Request){
	//User role name of the role to be assigned to user to get the role_id
	role := role.RoleModel{}
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to parse json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connect	
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT role_id FROM roles WHERE role_name = $1`
	row := db.QueryRow(stmt, role.RoleName)
	err := row.Scan(&role.RoleId)
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Get user_id from req params and convert it to int
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
	//Check if user belong to a group and if that group has that role to be assigned to the user
	stmt = `SELECT role_id FROM group_roles WHERE group_id =
	(SELECT group_id FROM users WHERE user_id = $1)`
	/*
		If above query returns null, it means either user does not belong to any group
		or the group the user belong to does not have that role
	*/
	var myUserId uint64
	row = db.QueryRow(stmt, userId)
	err = row.Scan(&myUserId)
	//Check if any row was returned or not
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error:   true,
			Message: "Incomplete!!! User is not part of a group or user's group does not have this role",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check for other errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	
	//Now update role_id of user on the users table
	stmt = `UPDATE users SET role_id = $2 WHERE user_id = $1`
	result, err := db.Exec(stmt, userId, role.RoleId)
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if any row was affected in the update operation
	if count, err := result.RowsAffected(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: "Internal server error" + err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}else{
		if count == 0 {
			w.WriteHeader(http.StatusNotModified)
			res := middleware.Response{
				Error: true,
				Message: "Unsuccessful!!! update operation",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything was fine, return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error: false,
		Message: "User role modified successfully",
	}
	json.NewEncoder(w).Encode(res)
}

//This is an update operation that sets user's role_id to null
 func RemoveRoleFromUser(w http.ResponseWriter, r *http.Request){
	 //Get user_id from req params and convert it to string
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
	 //Call db connection
	 db := db.ConnectDB()
	 defer db.Close()
	 stmt := `UPDATE users SET role_id = NULL WHERE user_id = $1`
	 result, err := db.Exec(stmt, userId)
	 // Check for errors
	 if err, ok := err.(*pq.Error); ok {
		 w.WriteHeader(http.StatusInternalServerError)
		 res := middleware.Response{
			 Error: true,
			 Message: "Internal server error" + err.Error(),
		 }
		 json.NewEncoder(w).Encode(res)
	 }
	 //Check if any row was affected during update
	 if count, err := result.RowsAffected(); err != nil {
		 w.WriteHeader(http.StatusInternalServerError)
		 res := middleware.Response{
			 Error: true,
			 Message: "Unable to return rows affected by update operation",
		 }
		 json.NewEncoder(w).Encode(res)
	 }else {
		 if count == 0 {
			 w.WriteHeader(http.StatusNotModified)
			 res := middleware.Response{
				 Error: true,
				 Message: "Unsuccessful!!! update operation" + err.Error(),
			 }
			 json.NewEncoder(w).Encode(res)
		 }
	 }
	 //If everything was fine, return response
	 w.WriteHeader(http.StatusCreated)
	 res := middleware.Response{
		 Error: true,
		 Message: "User role removed. User has no role at the moment. You need to assign one",
	 }
	 json.NewEncoder(w).Encode(res)
 }


