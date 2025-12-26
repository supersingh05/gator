package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/supersingh05/gator/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No Args Provided")
	}
	var userName sql.NullString
	userName.Scan(cmd.args[0])
	user := database.CreateUserParams{
		Name:      userName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
	}
	u, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}

	fmt.Printf("User Created: %v\n", u)
	err = s.config.SetUser(u.Name.String, u.ID)
	if err != nil {
		return err
	}
	return nil
}
