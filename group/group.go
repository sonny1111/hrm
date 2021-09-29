package group

import (
	"database/sql"
	"encoding/json"
	"hrm/db"
	"hrm/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func AddNewGroup(w http.ResponseWriter, r *http.Request) {
	group := GroupModel{}
	//Parse req body to json
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		res := middleware.Response{
			Error:   true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db conncetion
	db := db.ConnectDB()
	defer db.Close()
	stmt := `INSERT INTO groups(group_name, description) VALUES ($1, $2)`
	_, err = db.Exec(stmt)
	//Checking for errors
	if err, ok := err.(*pq.Error); ok {
		//Check for duplicate data
		if err.Code == "42701" || err.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			res := middleware.Response{
				Error:   true,
				Message: "Group already exists",
			}
			json.NewEncoder(w).Encode(res)
		}
	} else {
		//Checking for other errors
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//If everything went fine, return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error:   false,
		Message: "Group registered",
	}
	json.NewEncoder(w).Encode(res)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	//Convert req params to int
	params := mux.Vars(r)
	groupId, err := strconv.Atoi(params["group_id"])
	if err != nil {
		res := middleware.Response{
			Error:   true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	//Check if group has been assigned to user
	grpStmt := `SELECT user_id FROM users WHERE group_id = $1 LIMIT 1`
	row := db.QueryRow(grpStmt, groupId)
	var userId uint64

	err = row.Scan(&userId)
	//Checking for any error
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if group has roles or it's assigned users or both now
	if userId != 0 {
		w.WriteHeader(http.StatusFound)
		res := middleware.Response{
			Error: true,
			Message: "Cannot delete group that has been assigned to user!",
		}
		json.NewEncoder(w).Encode(res)
	}
	//You can now delete group
	stmt := `DELETE FROM groups WHERE id = $1`
	result, err := db.Exec(stmt, groupId)
	//Checking for errors
	if err, ok := err.(*pq.Error); ok {
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if any row was affected by the delete operation
	if count, err := result.RowsAffected(); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		res := middleware.Response{
			Error:   true,
			Message: "Unable to return rows affected by the delete operation",
		}
		json.NewEncoder(w).Encode(res)
	} else {
		if count == 0 {
			w.WriteHeader(http.StatusNotFound)
			res := middleware.Response{
				Error:   true,
				Message: "Delete operation NOT!!! successful. Group probably doesnt exist.",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything went fine, return response
	w.WriteHeader(http.StatusOK)
	res := middleware.Response{
		Error:   false,
		Message: "Group deleted successfully",
	}
	json.NewEncoder(w).Encode(res)
}

func GetGroup(w http.ResponseWriter, r *http.Request) {
	group := GroupModel{}
	//Convert req params to int
	params := mux.Vars(r)
	groupId, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   false,
			Message: "Could not convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT * FROM groups WHERE id = $1`
	row := db.QueryRow(stmt, groupId)
	err = row.Scan(stmt, &group.GroupId, &group.GroupName, &group.Description)
	//Check for no data found error
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error:   true,
			Message: "Group not found",
		}
		json.NewEncoder(w).Encode(res)
	} else {
		//For all other errors
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//If everything was fine, return group object
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(group)
}

//For fetching all the groups in the database
func GetGroups(w http.ResponseWriter, r *http.Request) {
	data := []GroupModel{}
	//Call db connecton
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT * FROM groups`
	rows, err := db.Query(stmt)
	//Checking for errors
	if err, ok := err.(*pq.Error); ok {
		//Check for no data found
		if err.Code == "P0002" || err.Code == "02000" {
			w.WriteHeader(http.StatusNotFound)
			res := middleware.Response{
				Error:   true,
				Message: "Group schema not yet populated",
			}
			json.NewEncoder(w).Encode(res)
		} else {
			//For all other errors
			w.WriteHeader(http.StatusInternalServerError)
			res := middleware.Response{
				Error:   true,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	defer rows.Close()
	for rows.Next() {
		group := GroupModel{}
		err := rows.Scan(&group.GroupId, &group.GroupName, &group.Description)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			res := middleware.Response{
				Error:   true,
				Message: "Unable to scan group result sets",
			}
			json.NewEncoder(w).Encode(res)
		}
		data = append(data, group)
	}
	//If no error, return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func EditGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	group := GroupModel{}
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Extract group_id from req params
	params := mux.Vars(r)
	groupId, err := strconv.Atoi(params["group_id"])
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
	stmt := `UPDATE groups SET group_name = $2 WHERE group_id = $1`
	result, err := db.Exec(stmt, uint64(groupId))
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		res := middleware.Response{
			Error: true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if any row was affected in the update operation
	if count, err := result.RowsAffected(); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		res := middleware.Response{
			Error: true,
			Message: "Error returning rows affected in the update operation",
		}
		json.NewEncoder(w).Encode(res)
	}else{
		if count == 0 {
			w.WriteHeader(http.StatusNotModified)
			res := middleware.Response{
				Error: true,
				Message: "Role probably doesnt exist",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything was fine, return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error: false,
		Message: "Group name modified successfully",
	}
	json.NewEncoder(w).Encode(res)
}
