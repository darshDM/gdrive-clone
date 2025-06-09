package types

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type StorageResponse struct {
	RemainingStorage int64 `json:"remaining_storage"`
}

type FileCreatedResponse struct {
	Message string `json:"message"`
}

type FileInfo struct {
	Filename string `json:"filename"`
	FileSize int64  `json:"file_size"`
}

type GetFilesResponse struct {
	Files []FileInfo `json:"files"`
}
