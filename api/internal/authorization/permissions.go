package authorization

const (
	PermissionCreateUser = "users:create"
	PermissionDeleteUser = "users:delete"
	PermissionViewUsers  = "users:view"
	PermissionUpdateUser = "users:update"

	PermissionCreateProduct = "products:create"
	PermissionUpdateProduct = "products:update"
	PermissionDeleteProduct = "products:delete"

	PermissionCreateOrder   = "orders:create"
	PermissionViewOrder     = "orders:view"
	PermissionViewAllOrders = "orders:viewall"
	PermissionCancelOrder   = "orders:cancel"
)

const (
	RoleAdmin = "Admin"
	RoleUser  = "User"
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
		PermissionViewAllOrders,
		PermissionCancelOrder,
	},
	RoleUser: {
		PermissionUpdateUser,
		PermissionCreateOrder,
		PermissionViewOrder,
		PermissionCancelOrder,
	},
}
