package data

import (
	"database/sql"
	"time"
)

type PostgresTestRepository struct {
	conn *sql.DB
}

func NewPostgresTestRepository(conn *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{conn}
}

// All returns a slice of all users, sorted by lastname
func (p *PostgresTestRepository) All() ([]*User, error) {
	users := []*User{}

	return users, nil
}

// ByEmail returns user by email
func (p *PostgresTestRepository) ByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		FirstName: "Joseph",
		LastName:  "Gilgun",
		Email:     "email",
		Password:  "password",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// ByID returns one user by id
func (p *PostgresTestRepository) ByID(id int) (*User, error) {
	user := User{
		ID:        1,
		FirstName: "Joseph",
		LastName:  "Gilgun",
		Email:     "email",
		Password:  "password",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (p *PostgresTestRepository) Update(user User) error {
	return nil
}

// Delete deletes one user from the database, by User.ID
func (p *PostgresTestRepository) Delete(id int) error {
	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (p *PostgresTestRepository) Insert(user User) (int, error) {
	return 2, nil
}

// ResetPassword is the method we will use to change a user's password.
func (p *PostgresTestRepository) ResetPassword(password string, user User) error {
	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (p *PostgresTestRepository) PasswordMatches(plainText string, user User) (bool, error) {
	return true, nil
}
