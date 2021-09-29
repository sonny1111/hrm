package middleware

import (
	"encoding/json"
	"hrm/db"
	"net/http"

	"github.com/lib/pq"
)
type privModel struct {
	PrivilegeName string 
}

func IsAuthorize(allowedPrivilege string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := []privModel{}
		//Extract Roleid from the context of the jwt middleware
		roleId, ok := r.Context().Value("role_id").(int)
		if !ok {
			res := Response{
				Error:   true,
				Message: "Unable to extract permission info",
			}
			json.NewEncoder(w).Encode(res)
		}
		//Check the database for privileges assigned to this role
		db := db.ConnectDB()
		defer db.Close()
		stmt := `SELECT privilege_name FROM privileges
					WHERE privilege_name
					IN
		SELECT privilege_id FROM role_privileges WHERE role_id = $1`
		rows, err := db.Query(stmt, uint64(roleId))
		
		//Check for all errors
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "P0002" || err.Code == "02000" {
				w.WriteHeader(http.StatusNotFound)
				res := Response{
					Error:   true,
					Message: "No permission info found",
				}
				json.NewEncoder(w).Encode(res)
			}else {
				w.WriteHeader(http.StatusInternalServerError)
				res := Response{
					Error: true,
					Message: "Internal server error" + err.Error(),
				}
				json.NewEncoder(w).Encode(res)
			}
		}
		defer rows.Close()
		for rows.Next() {
			pName := privModel{}
			err := rows.Scan(&pName.PrivilegeName)
			if err != nil {
				res := Response{
					Error:   true,
					Message: "Unable to scan permission variables",
				}
				json.NewEncoder(w).Encode(res)
			}
			data = append(data, pName)
			//Convert above data struct to a slice of string
			var priviliges []string
			for _, v := range data {
				priviliges = append(priviliges, v.PrivilegeName)
			}
			//Check if privileges slice contain privilege allowed for the this endpoint
			condition := contains(priviliges, allowedPrivilege)
			if !condition  {
				res := Response{
					Error:   true,
					Message: "Unauthorized",
				}
				json.NewEncoder(w).Encode(res)
			}

		}
		next.ServeHTTP(w, r)

	})
}
