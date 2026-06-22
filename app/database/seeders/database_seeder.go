package seeders

import "log"

// SeedDatabase seeds the database
func SeedDatabase() error {
	admin, err := SeedAdmins()
	if err != nil {
		return err
	}
	if admin != nil {
		log.Printf("Seeded %s Super Admin\n", admin.Name)
	}

	manager, err := SeedManagers()
	if err != nil {
		return err
	}
	if manager != nil {
		log.Printf("Seeded %s Agency Manager\n", manager.Name)
	}

	return nil
}
