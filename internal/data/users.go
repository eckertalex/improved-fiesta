package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eckertalex/improved-fiesta/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

const (
	UserRole  = "user"
	AdminRole = "admin"
)

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Role      string    `json:"role"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (u *User) IsAdmin() bool {
	return u.Role == AdminRole
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateRole(v *validator.Validator, role string) {
	v.Check(role != "", "role", "must be provided")
	v.Check(role == UserRole || role == AdminRole, "role", "must be either 'user' or 'admin'")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "name", "must be provided")
	v.Check(len(user.Username) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)
	ValidateRole(v, user.Role)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) GetAll(username string, email string, filter Filters) ([]*User, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, updated_at, username, email, password_hash, activated, role, version
		FROM users
		WHERE (username LIKE ? OR ? = '')
		AND (email LIKE ? OR ? = '')
		ORDER BY %s %s, id ASC
		LIMIT ? OFFSET ?
	`, filter.sortColumn(), filter.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	usernamePattern := "%" + username + "%"
	emailPattern := "%" + email + "%"
	args := []any{usernamePattern, username, emailPattern, email, filter.limit(), filter.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Username,
			&user.Email,
			&user.Password.hash,
			&user.Activated,
			&user.Role,
			&user.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return users, metadata, nil
}

func (m UserModel) Insert(user *User) error {
	if user.Role == "" {
		user.Role = UserRole
	}

	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)

	query := `
		INSERT INTO users (username, email, password_hash, activated, role)
		VALUES (?, ?, ?, ?, ?)
	`

	args := []any{user.Username, user.Email, user.Password.hash, user.Activated, user.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return ErrDuplicateEmail
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return ErrDuplicateUsername
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fetchedUser, err := m.GetByID(id)
	if err != nil {
		return err
	}

	user.ID = id
	user.CreatedAt = fetchedUser.CreatedAt
	user.UpdatedAt = fetchedUser.UpdatedAt
	user.Version = fetchedUser.Version

	return nil
}

func (m UserModel) GetByID(id int64) (*User, error) {
	query := `
		SELECT id, created_at, updated_at, username, email, password_hash, activated, role, version
		FROM users
		WHERE id = ?
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Role,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	email = strings.ToLower(email)

	query := `
		SELECT id, created_at, updated_at, username, email, password_hash, activated, role, version
		FROM users
		WHERE email = ?
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Role,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)

	query := `
		UPDATE users
		SET username = ?, email = ?, password_hash = ?, activated = ?, role = ?, version = version + 1
		WHERE id = ? AND version = ?
	`

	args := []any{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.Role,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return ErrDuplicateEmail
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return ErrDuplicateUsername
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrEditConflict
	}

	user.Version++

	return nil
}

func (m UserModel) CountAdminUsers() (int, error) {
	query := `
        SELECT COUNT(*) 
        FROM users 
        WHERE role = ?
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := m.DB.QueryRowContext(ctx, query, AdminRole).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
		SELECT users.id, users.created_at, users.updated_at, users.username, users.email, users.password_hash, users.activated, users.role, users.version
		FROM users
		INNER JOIN tokens
		ON users.id = tokens.user_id
		WHERE tokens.hash = ?
		AND tokens.scope = ?
		AND tokens.expiry > ?
	`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Role,
		&user.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (m UserModel) Delete(id int64) error {
	user, err := m.GetByID(id)
	if err != nil {
		return err
	}

	query := `
        DELETE FROM users
        WHERE id = ? AND version = ?
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id, user.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrEditConflict
	}

	return nil
}
