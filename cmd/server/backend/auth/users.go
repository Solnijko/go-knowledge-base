package auth

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Random password generation
func generatePassword(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// User creation function
func CreateUser(pool *pgxpool.Pool, ctx context.Context, user User, logger *slog.Logger) error {
	// Random password is generated if it was not set
	if user.Password == "" {
		logger.Info("generating random temporary password since it was not set")
		user.Password = generatePassword(12)
		// TODO: Implement password sending after SMTP implementation
		// Example: SendWelcomeEmail(email, password)
	}

	query := `INSERT INTO users (username, email, password, access, created) VALUES ($1, $2, $3, $4, $5);`
	commandTag, err := pool.Exec(ctx, query, user.Username, user.Email, user.Password, user.Access, user.Created)
	logger.Debug("executing user creation query", "command_tag", commandTag)

	if err != nil {
		logger.Error("unable to create user", "err", err)
		return err
	}
	logger.Info("created user", "username", user.Username, "email", user.Email, "access", user.Access)
	return nil
}

// Credentials must be set in environment
func FirstUser(email string, username string, logger *slog.Logger) (User, error) {
	var firstUser User

	userPassword := generatePassword(12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("password generation error", "err", err)
		return User{}, err
	}

	firstUser.Username = username
	firstUser.Email = email
	firstUser.Password = string(hashedPassword)
	firstUser.Access = []string{"admin"}
	firstUser.Created = time.Now()

	// TODO: Implement password sending after SMTP implementation
	// Example: SendWelcomeEmail(email, password)
	logger.Info("created first user", "email", email, "username", username, "access", firstUser.Access)
	return firstUser, nil
}
