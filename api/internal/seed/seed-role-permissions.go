package seed

import (
	"context"
	"database/sql"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	"github.com/rawbil/ecom2/internal/authorization"
)

func SeedRolePermissions(ctx context.Context, repository repository.Queries, db *sql.DB) error {
	for roleName, permissions := range authorization.RolePermissions {
		role_id, err := repository.GetRoleID(ctx, roleName)
		if err != nil {
			return err
		}

		for _, permission := range permissions {
			permission_id, err := repository.GetPermissionID(ctx, permission)
			if err != nil {
				return err
			}

			query := `
			INSERT IGNORE INTO role_permissions(role_id, permission_id)
			VALUES (?, ?)
			`

			if _, err := db.Exec(query, role_id, permission_id); err != nil {
				return err
			}
		}
	}

	return nil
}
