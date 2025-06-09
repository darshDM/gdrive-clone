package store

import (
	"context"
	"fmt"
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	CreatedOn    string `json:"created_on"`
	TotalStorage int64  `json:"total_storage"`
	UsedStorage  int64  `json:"used_storage"`
}

func (driver *Store) CreateNewUser(ctx context.Context, user *User) error {
	statement, err := driver.Db.Prepare("INSERT INTO users (username, password, total_storage, used_storage) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.Username, user.Password, user.TotalStorage, user.UsedStorage)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

func (driver *Store) GetUser(ctx context.Context, username string) (*User, error) {
	statement, err := driver.Db.Prepare("SELECT id, username, password, created_at, total_storage, used_storage FROM users WHERE username = ?")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	row := statement.QueryRow(username)
	user := &User{}
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedOn, &user.TotalStorage, &user.UsedStorage)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	return user, nil
}
