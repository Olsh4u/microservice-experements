package data

type Repository interface {
	All() ([]*User, error)
	ByEmail(email string) (*User, error)
	ByID(id int) (*User, error)
	Update(user User) error
	Delete(id int) error
	Insert(user User) (int, error)
	ResetPassword(password string, user User) error
	PasswordMatches(plainText string, user User) (bool, error)
}
