package seed

import (
	"database/sql"

	"github.com/rawbil/ecom2/internal/authorization"
)

func SeedPermissions(db *sql.DB) error {
	permissions := []string{
		authorization.PermissionCreateUser,
		authorization.PermissionDeleteUser,
		authorization.PermissionViewUsers,
		authorization.PermissionUpdateUser,
		authorization.PermissionCreateProduct,
		authorization.PermissionUpdateProduct,
		authorization.PermissionDeleteProduct,
		authorization.PermissionCreateOrder,
		authorization.PermissionViewOrder,
		authorization.PermissionViewAllOrders,
		authorization.PermissionCancelOrder,
	}

	for _, permission := range permissions {
		query := `
		INSERT IGNORE INTO user_permissions(permission)
		VALUES (?)
		`

		_, err := db.Exec(query, permission)
		if err != nil {
			return err
		}
	}

	return nil
}
