package usecases

import (
	"context"

	"github.com/vorotilkin/file-storage/domain/models"
	"github.com/vorotilkin/file-storage/proto"
)

type FilesRepository interface {
	Create(ctx context.Context, file models.CreateFileRequest) (models.CreateFileResponse, error)
	ObjectKeys(ctx context.Context, fileIDs []int32) (map[int32]string, error)
}

type S3Service interface {
	PresignPut(ctx context.Context, objectName string) (string, error)
	PresignGet(ctx context.Context, objectName string) (string, error)
	Bucket() string
}

type FileStorageServer struct {
	proto.UnimplementedFileStorageServiceServer

	filesRepo FilesRepository
	s3Service S3Service
}

func (s *FileStorageServer) RegisterFile(ctx context.Context, request *proto.RegisterFileRequest) (*proto.RegisterFileResponse, error) {
	filename := request.GetFilename()
	entityName := request.GetEntityName()
	objectKey := models.CreateObjectKey(filename, entityName).String()

	createdFile, err := s.filesRepo.Create(ctx, models.CreateFileRequest{
		Filename:    filename,
		ContentType: request.GetContentType(),
		Bucket:      s.s3Service.Bucket(),
		ObjectKey:   objectKey,
	})
	if err != nil {
		return nil, err
	}

	url, err := s.s3Service.PresignPut(ctx, objectKey)
	if err != nil {
		return nil, err
	}

	return &proto.RegisterFileResponse{
		FileId: createdFile.ID,
		PutUrl: url,
	}, nil
}

func (s *FileStorageServer) DownloadLink(ctx context.Context, request *proto.DownloadLinkRequest) (*proto.DownloadLinkResponse, error) {
	objectKeyMap, err := s.filesRepo.ObjectKeys(ctx, request.GetFileIds())
	if err != nil {
		return nil, err
	}

	fileURLsMap := make(map[int32]string, len(objectKeyMap))

	for fileID, objectKey := range objectKeyMap {
		url, err := s.s3Service.PresignGet(ctx, objectKey)
		if err != nil {
			return nil, err
		}

		fileURLsMap[fileID] = url
	}

	return &proto.DownloadLinkResponse{
		FileUrlsMap: fileURLsMap,
	}, nil
}

func NewFileStorageServer(
	filesRepo FilesRepository,
	s3Service S3Service,
) *FileStorageServer {
	return &FileStorageServer{
		filesRepo: filesRepo,
		s3Service: s3Service,
	}
}
