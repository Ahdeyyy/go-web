package database

import (
	"github.com/Ahdeyyy/go-web/internal/models"
)

// dbInterface is the interface for the database
type DBInterface interface {
	AllUsers() ([]models.User, error)
	GetUserByID(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	CreateUser(u models.User) (int, error)
	UpdateUser(u models.User) error
	DeleteUser(id int) error
	Authenticate(email, password string) (models.User, string, error)
}
