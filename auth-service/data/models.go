package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = 3 * time.Second

type PostgresRepository struct {
	Conn *sql.DB
}

func NewPostgresRepository(conn *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		Conn: conn,
	}
}

// User is the structure which holds one user from database
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// All returns a slice of all users, sorted by lastname
func (p *PostgresRepository) All() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
				from users order by last_name`

	rows, err := p.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// ByEmail returns user by email
func (p *PostgresRepository) ByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at 
				from users where email = $1`

	var user User
	row := p.Conn.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("Error while scanning row:", err)
		return nil, err
	}

	return &user, nil
}

// ByID returns one user by id
func (p *PostgresRepository) ByID(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
				from users where id = $1`

	var user User
	row := p.Conn.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("Error while scanning row:", err)
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (p *PostgresRepository) Update(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
				email = $1,
				first_name = $2,
				last_name = $3,
				user_active = $4,
				updated_at = $5
				where id = $6`

	_, err := p.Conn.ExecContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		time.Now(),
		user.ID,
	)

	if err != nil {
		log.Println("Error while updating user:", err)
		return err
	}

	return nil
}

// Delete deletes one user from the database, by User.ID
func (p *PostgresRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := p.Conn.ExecContext(ctx, stmt, id)
	if err != nil {
		log.Println("Error while deleting user:", err)
		return err
	}

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (p *PostgresRepository) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = p.Conn.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPass,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		log.Println("Error while inserting user:", err)
		return 0, err
	}

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (p *PostgresRepository) ResetPassword(password string, user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = p.Conn.ExecContext(ctx, stmt, hashedPass, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (p *PostgresRepository) PasswordMatches(plainText string, user User) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			//invalid pass
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
