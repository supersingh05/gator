package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No Args Provided")
	}
	user := cmd.args[0]
	var userName sql.NullString
	userName.Scan(user)
	if !userName.Valid {
		return errors.New("couldnt parse sql nullstring")
	}

	u, err := s.db.GetUserByName(context.Background(), userName)
	if err != nil {
		return err
	}
	err = s.config.SetUser(u.Name.String, u.ID)
	if err != nil {
		return err
	}
	fmt.Printf("User %s has been set", user)
	return nil
}
