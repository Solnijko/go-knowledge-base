package auth

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Random password generation
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*"

func generatePassword(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Credentials must be set in environment
func FirstUser() (User, error) {
	var firstUser User

	email := os.Getenv("GOKB_EMAIL")
	username := os.Getenv("GOKB_USERNAME")

	if username == "" {
		username = "gokb"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(generatePassword(12)), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("password generation error %v", err)
	}

	if email == "" {
		return User{}, fmt.Errorf("email was not set")
	}

	firstUser.Username = username
	firstUser.Email = email
	firstUser.Password = string(hashedPassword)
	firstUser.Access = []string{"admin"}
	firstUser.Created = time.Now()

	log.Printf("Generated first user with email: %s, username: %s", email, username)
	return firstUser, nil
}

func CreateUserTable(pool *pgxpool.Pool, ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password VARCHAR(60) NOT NULL,
		access TEXT[] NOT NULL,
		created TIMESTAMP NOT NULL
	);`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to create table: %w", err)
	}
	return nil
}

func CreateUser(pool *pgxpool.Pool, ctx context.Context, user User) error {
	query := `INSERT INTO users (username, email, password, access, created) VALUES ($1, $2, $3, $4, $5);`
	_, err := pool.Exec(ctx, query, user.Username, user.Email, user.Password, user.Access, user.Created)
	if err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}
	return nil
}
