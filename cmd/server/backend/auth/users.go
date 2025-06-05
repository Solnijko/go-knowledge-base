package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Random password generation
func GeneratePassword(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*"
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random password: %v", err)
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}

// User creation function
func CreateUser(pool *pgxpool.Pool, ctx context.Context, user User, logger *slog.Logger) error {
	// Check if user already exists
	var userCount int
	checkQuery := `SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2`
	err := pool.QueryRow(ctx, checkQuery, user.Username, user.Email).Scan(&userCount)
	logger.Debug("executing user check query", "query", checkQuery)
	if err != nil {
		logger.Error("unable to count users", "err", err)
		return err
	}
	if userCount != 0 {
		logger.Warn("user with this username or email already exists", "username", user.Username, "email", user.Email)
		return fmt.Errorf("user with username %s or email %s already exists", user.Username, user.Email)
	}

	// Random password is generated if it was not set
	if user.Password == "" {
		logger.Error("password was not set")
		return fmt.Errorf("password was not set")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("password hashing error", "err", err)
		return err
	}

	if len(user.Access) == 0 {
		return fmt.Errorf("user access is empty")
	}

	creationQuery := `INSERT INTO users (username, email, password, access, created) VALUES ($1, $2, $3, $4, $5);`
	_, err = pool.Exec(ctx, creationQuery, user.Username, user.Email, hashedPassword, user.Access, user.Created)
	logger.Debug("executing user creation query", "query", creationQuery)
	if err != nil {
		logger.Error("unable to create user", "err", err)
		return err
	}

	logger.Info("created user", "username", user.Username, "email", user.Email, "access", user.Access)
	return nil
}

// Credentials must be set in environment
func FirstUser(email string, username string, password string, logger *slog.Logger, pool *pgxpool.Pool) (User, error) {
	var firstUser User
	firstUser.Username = username
	firstUser.Email = email
	firstUser.Password = password
	firstUser.Created = time.Now()
	// TODO: Implement password sending after SMTP implementation
	// Example: SendWelcomeEmail(email, password)
	fmt.Println(firstUser.Password)
	firstUser.Access = []string{"admin"}

	err := CreateUser(pool, context.Background(), firstUser, logger)
	if err != nil {
		logger.Error("failed to create first user", "err", err)
		return User{}, err
	}

	logger.Info("created first user", "email", email, "username", username, "access", firstUser.Access)
	return firstUser, nil
}
