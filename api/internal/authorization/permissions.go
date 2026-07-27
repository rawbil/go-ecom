package authorization

const (
	PermissionCreateUser = "users:create"
	PermissionDeleteUser = "users:delete"
	PermissionViewUsers  = "users:view"
	PermissionUpdateUser = "users:update"

	PermissionCreateProduct = "products:create"
	PermissionUpdateProduct = "products:update"
	PermissionDeleteProduct = "products:delete"

	PermissionCreateOrder = "orders:create"
	PermissionViewOrder   = "orders:view"
	PermissionCancelOrder = "orders:cancel"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

var RolePermissions = map[string][]string{
	RoleAdmin: {
		PermissionCreateUser,
		PermissionDeleteUser,
		PermissionViewUsers,
		PermissionUpdateUser,
		PermissionCreateProduct,
		PermissionUpdateProduct,
		PermissionDeleteProduct,
		PermissionCreateOrder,
		PermissionViewOrder,
		PermissionCancelOrder,
	},
	RoleUser: {
		PermissionUpdateUser,
		PermissionCreateOrder,
		PermissionViewOrder,
		PermissionCancelOrder,
	},
}
