package storage

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/darshDM/gdrive-clone-api/internal/store"
	"github.com/darshDM/gdrive-clone-api/types"
)

type StorageService struct {
	store        store.Store
	UploadFolder string
}

func NewStorageService(s store.Store) *StorageService {
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Fatal("Error while initializing storage service", err.Error())
	}
	return &StorageService{
		store:        s,
		UploadFolder: "uploads",
	}
}

func (s *StorageService) UploadFile(ctx context.Context, file multipart.File, fileInfo *multipart.FileHeader) error {
	user, ok := ctx.Value("user").(*store.User)
	if !ok {
		return errors.New("user not found in context")
	}
	remainingStorage := user.TotalStorage - user.UsedStorage
	if remainingStorage < fileInfo.Size {
		return errors.New("not enough storage space")
	}
	fileName := filepath.Clean(fileInfo.Filename)
	filePath := filepath.Join(s.UploadFolder, user.Username, fileName)
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return errors.New("failed to create directory")
	}

	storedFile, err := os.Create(filePath)
	if err != nil {
		return errors.New("failed to create file")
	}
	defer storedFile.Close()

	if _, err := io.Copy(storedFile, file); err != nil {
		return errors.New("failed to copy file")
	}
	if err := s.store.UpdateStorage(ctx, user, fileInfo.Size); err != nil {
		return errors.New("failed to update storage")
	}
	return nil
}

func (storageService *StorageService) GetRemainingStorage(ctx context.Context) (int64, error) {
	user, ok := ctx.Value("user").(*store.User)
	if !ok {
		return 0, errors.New("user not found in context")
	}
	remainingStorage := user.TotalStorage - user.UsedStorage
	return remainingStorage, nil

}

func (s *StorageService) GetFiles(ctx context.Context, limit int, offset int) (types.GetFilesResponse, error) {
	user, ok := ctx.Value("user").(*store.User)
	if !ok {
		return types.GetFilesResponse{}, errors.New("user not found in context")
	}

	userFilesPath := filepath.Join(s.UploadFolder, user.Username)
	entries, err := os.ReadDir(userFilesPath)
	if err != nil {
		return types.GetFilesResponse{}, errors.New("failed to read user files directory")
	}
	var filesInfo types.GetFilesResponse
	for i, file := range entries {
		if i < offset {
			continue
		}
		if i >= offset+limit {
			break
		}
		f, err := getFileInfo(file)
		if err != nil {
			return types.GetFilesResponse{}, errors.New("failed to get file info")
		}
		filesInfo.Files = append(filesInfo.Files, f)
	}
	return filesInfo, nil
}

func getFileInfo(file fs.DirEntry) (types.FileInfo, error) {
	fileInfo, err := file.Info()
	if err != nil {
		return types.FileInfo{}, err
	}
	return types.FileInfo{
		Filename: fileInfo.Name(),
		FileSize: fileInfo.Size(),
	}, nil
}
