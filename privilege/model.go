package privilege


type PrivilegeModel struct{
	PrivilegeId uint64 `json:"id"`
	PrivilegeName string `json:"privilege_name"`
	Description string `json:"description"`
}