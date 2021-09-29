package group

import (
	"database/sql"
	"encoding/json"
	"hrm/db"
	"hrm/middleware"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

//This is an update on user table. No need to check whether the user has the group assigned already or not
func AddUserToGroup(w http.ResponseWriter, r *http.Request){
	//Use group name to extact the group id of the group to be assigned to user
	group := GroupModel{}
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil{
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT group_id FROM groups WHERE group_name = $1`
	row := db.QueryRow(stmt, group.GroupName)
	//Check if any row is returned or not after scanning the return result
	err := row.Scan(&group.GroupId)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error: true,
			Message: "Group not found",
		}
		json.NewEncoder(w).Encode(res)
	}
	
//  user := user.UserModel{}
 //Extract user_id from request body
  params := mux.Vars(r)
  //Convert user id to int
   userId, err := strconv.Atoi(params["user_id"])
   if err != nil {
	   w.WriteHeader(http.StatusBadRequest)
	   res := middleware.Response{
		   Error: true,
		   Message: "Unable to convert req params to int",
	   }
	   json.NewEncoder(w).Encode(res)
   }
   //Now update user with the new group id
   stmt = `UPDATE users SET group_id = $2 WHERE user_id = $1`
   result, err := db.Exec(stmt, userId, group.GroupId)
   //Check for errors
   if err, ok := err.(*pq.Error); ok {
	   w.WriteHeader(http.StatusInternalServerError)
	   res := middleware.Response{
		   Error: true,
		   Message: "Internal server error" + err.Error(),
	   }
	   json.NewEncoder(w).Encode(res)
   }
   //Check if any row affected
  if count, err := result.RowsAffected(); err != nil {
	  w.WriteHeader(http.StatusExpectationFailed)
	  res := middleware.Response{
		  Error: true,
		  Message: "Unable to count rows affected in update operation",
	  }
	  json.NewEncoder(w).Encode(res)
  } else{
	  if count == 0 {
		  w.WriteHeader(http.StatusNotModified)
		  res := middleware.Response{
			  Error: true,
			  Message: "Update operation NOT!!! successful. User probably dont exist",
		  }
		  json.NewEncoder(w).Encode(res)
	  }
  }
  //If update operation was successful, return response
  w.WriteHeader(http.StatusCreated)
  res := middleware.Response{
	  Error: false,
	  Message: "User group object modified.",
  }
  json.NewEncoder(w).Encode(res)
}

//this is an update operation that sets group_id on users table to null
func RemoveUserFromGroup(w http.ResponseWriter, r *http.Request){
	//Extract user_d from req params
	params := mux.Vars(r)
	//Convert req params to int
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
	stmt := `UPDATE users SET group_id = NULL WHERE user_id = $1`
	result, err := db.Exec(stmt, userId)
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
	if count, err := result.RowsAffected(); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		res := middleware.Response{
			Error: true,
			Message: "Unable to count rows affected by the update operation",
		}
		json.NewEncoder(w).Encode(res)
	}else{
		if count == 0 {
			w.WriteHeader(http.StatusNotModified)
			res := middleware.Response{
				Error: true,
				Message: "Update operation NOT!!! successful. User probably doesnt exist.",
			}
			json.NewEncoder(w).Encode(res)
		}
		//If everything went fine, return response
		w.WriteHeader(http.StatusCreated)
		res := middleware.Response{
			Error: false,
			Message: "User removed from group successfully",
		}
		json.NewEncoder(w).Encode(res)
	}
}

func RemoveAllUsersFromGroup(w http.ResponseWriter, r *http.Request){
	//Extract group id from req params and convert to int
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
	stmt := `DELETE FROM users WHERE group_id = $1`
	result, err := db.Exec(stmt, groupId)
	//Checking for errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: "Internal server error" + err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	if count, err := result.RowsAffected(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: "Error counting rows returned by the delete operation.",
		}
		json.NewEncoder(w).Encode(res)
	}else {
		if count == 0 {
			w.WriteHeader(http.StatusNotModified)
			res := middleware.Response{
				Error: true,
				Message: "No row affected by the delete operation!!!",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If operation was successful, return response
	w.WriteHeader(http.StatusOK)
	res := middleware.Response{
		Error: false,
		Message: "Operation successful",
	}
	json.NewEncoder(w).Encode(res)
	
}