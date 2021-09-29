package role


type RoleModel struct{
	RoleId uint64 `json:"id"`
	RoleName string `json:"role_name"`
	Description string `json:"description"`
	
}