package store

import "context"

type FileInfo struct {
	Name string
	Size int64
}

type FileArray struct {
	Files []FileInfo
}

func (s *Store) UpdateStorage(ctx context.Context, user *User, size int64) error {
	statement, err := s.Db.Prepare("UPDATE users SET used_storage = used_storage + ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(size, user.ID)
	if err != nil {
		return err
	}
	return nil
}
