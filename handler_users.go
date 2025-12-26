package main

import (
	"context"
	"fmt"
)

func handleUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	name := s.config.CurrentUserName
	for _, u := range users {
		if u.Name.String == name {
			fmt.Printf("* %s (current)\n", u.Name.String)
		} else {
			fmt.Printf("* %s\n", u.Name.String)
		}
	}
	return nil
}
