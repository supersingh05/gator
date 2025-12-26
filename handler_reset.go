package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		fmt.Println("failed to reset the users table")
		return err
	}
	fmt.Println("reset the users table")
	err = s.db.DeleteAllFeeds(context.Background())
	if err != nil {
		fmt.Println("failed to reset the users table")
		return err
	}
	fmt.Println("reset the feeds table")
	return nil
}
