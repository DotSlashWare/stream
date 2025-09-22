package profile

import (
	"log"

	"github.com/artumont/DotSlashStream/backend/internal/database/postgres"
)

type Service struct {
	Database *postgres.Manager
}

func NewService(db *postgres.Manager) *Service {
	return &Service{
		Database: db,
	}
}

// GetAllProfiles retrieves all profiles from the database.
func (service *Service) GetAllProfiles() ([]Profile, error) {
	dbManager := service.Database

	rows, err := dbManager.SelectFrom("profiles", []string{"id", "username", "is_child", "created_at", "updated_at", "last_login"}, "")
	if err != nil {
		log.Printf("Error fetching profiles: %v", err)
		return nil, err
	}

	var profiles []Profile
	for rows.Next() {
		var profile Profile
		if err := rows.Scan(&profile.Id, &profile.Username, &profile.IsChild, &profile.CreatedAt, &profile.UpdatedAt, &profile.LastLogin); err != nil {
			log.Printf("Error scanning profile row: %v", err)
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	
	return profiles, nil
}
