package files

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/samber/lo"
	"github.com/vorotilkin/file-storage/domain/models"
	"github.com/vorotilkin/file-storage/pkg/database"
	"github.com/vorotilkin/file-storage/schema/gen/file_storage/public/model"
	"github.com/vorotilkin/file-storage/schema/gen/file_storage/public/table"
)

type Repository struct {
	conn *database.Database
}

func (r *Repository) Create(ctx context.Context, file models.CreateFileRequest) (models.CreateFileResponse, error) {
	sql, args := table.Files.
		INSERT(
			table.Files.Filename,
			table.Files.ContentType,
			table.Files.Bucket,
			table.Files.ObjectKey,
		).
		MODEL(
			model.Files{
				Bucket:      file.Bucket,
				ObjectKey:   file.ObjectKey,
				Filename:    file.Filename,
				ContentType: file.ContentType,
				CreatedAt:   time.Now(),
			},
		).
		RETURNING(table.Files.ID).
		Sql()

	row := r.conn.QueryRow(ctx, sql, args...)
	dbFile := model.Files{}

	err := row.Scan(&dbFile.ID)
	if err != nil {
		return models.CreateFileResponse{}, err
	}

	return models.CreateFileResponse{
		ID: dbFile.ID,
	}, nil
}

func (r *Repository) Upload(ctx context.Context, fileID int32) (bool, error) {
	sql, args := table.Files.
		UPDATE(
			table.Files.UploadedAt,
		).
		MODEL(
			model.Files{
				UploadedAt: lo.ToPtr(time.Now()),
			},
		).
		WHERE(table.Files.ID.EQ(postgres.Int(int64(fileID)))).
		Sql()

	result, err := r.conn.Exec(ctx, sql, args...)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func (r *Repository) ObjectKey(ctx context.Context, fileID int32) (string, error) {
	sql, args := table.Files.
		SELECT(
			table.Files.ObjectKey,
		).
		WHERE(table.Files.ID.EQ(postgres.Int(int64(fileID)))).
		Sql()

	row := r.conn.QueryRow(ctx, sql, args...)
	dbFile := model.Files{}

	err := row.Scan(&dbFile.ObjectKey)
	if err != nil {
		return "", err
	}

	return dbFile.ObjectKey, nil
}

func NewRepository(conn *database.Database) *Repository {
	return &Repository{conn: conn}
}
