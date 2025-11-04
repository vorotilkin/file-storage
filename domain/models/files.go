package models

type CreateFileRequest struct {
	Filename    string
	ContentType string
	Bucket      string
	ObjectKey   string
}

type CreateFileResponse struct {
	ID int32
}
