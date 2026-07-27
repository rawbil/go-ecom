package seed

import "database/sql"

func SeedRoles(db *sql.DB) error {
	roles := []string{
		"Admin",
		"User",
	}

	for _, role := range roles {
		query := `INSERT IGNORE INTO user_roles (role) VALUES (?)`
		_, err := db.Exec(query, role)
		if err != nil {
			return err
		}
	}

	return nil
}
