package sqlite

import (
	"context"
	"errors"
	"time"

	"github.com/Ahdeyyy/go-web/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AllUsers returns all users from the database
func (m *sqliteDBinterface) AllUsers() ([]models.User, error) {
	rows, err := m.DB.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// CreateUser inserts a user into the database
func (m *sqliteDBinterface) CreateUser(u models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exists, _ := m.GetUserByEmail(u.Email)
	if exists.ID != 0 {
		return -1, errors.New("user already exists")
	}

	exists, _ = m.GetUserByUsername(u.Username)
	if exists.ID != 0 {
		return -1, errors.New("user already exists")
	}

	stmt := `INSERT INTO users (firstname, lastname, email, password, created_at, updated_at) VALUES(?,?,?,?,?,?)`
	password, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	result, err := m.DB.ExecContext(ctx, stmt, u.Firstname, u.Lastname, u.Email, password, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetUserByID returns a user by id
func (m *sqliteDBinterface) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT * FROM users WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	var u models.User
	err := row.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

// GetUserByEmail returns a user by email
func (m *sqliteDBinterface) GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT * FROM users WHERE email = ?`
	row := m.DB.QueryRowContext(ctx, stmt, email)

	var u models.User
	err := row.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

// GetUserByUsername returns a user by username
func (m *sqliteDBinterface) GetUserByUsername(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT * FROM users WHERE username = ?`
	row := m.DB.QueryRowContext(ctx, stmt, username)

	var u models.User
	err := row.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates a user in the database
func (m *sqliteDBinterface) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE users SET firstname = ?, lastname = ?, email = ?, password = ?, updated_at = ? WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, u.Firstname, u.Lastname, u.Email, u.Password, time.Now(), u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user from the database
func (m *sqliteDBinterface) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM users WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate checks if a user exists and the password is correct
func (m *sqliteDBinterface) Authenticate(email string, password string) (models.User, string, error) {
	u, err := m.GetUserByEmail(email)
	if err != nil {
		return models.User{}, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return models.User{}, "", errors.New("incorrect password")
	} else if err != nil {
		return models.User{}, "", err
	}

	return u, u.Password, nil
}
