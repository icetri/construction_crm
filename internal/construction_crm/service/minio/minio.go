package minio

import (
	"context"
	"fmt"
	"github.com/construction_crm/internal/construction_crm/types/config"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mime/multipart"
	"time"
)

type File struct {
	ID       string    `db:"id"`
	Date     time.Time `db:"created_at"`
	Tag      string    `db:"tag"`
	Url      string    `db:"url"`
	Length   int64     `db:"length"`
	MimeType string    `db:"mime"`
	Bucket   string    `db:"bucket"`
	Object   string    `db:"object"`
	Role     string    `db:"role"`
	UserId   int       `db:"user_id"`
}

type FileStorage struct {
	client *minio.Client
}

func NewMinio(cfg *config.Storage) (*FileStorage, error) {
	client, err := minio.New(
		cfg.Host,
		&minio.Options{Creds: credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""), Secure: cfg.SSL})

	if err != nil {
		return nil, err
	}

	return &FileStorage{
		client: client,
	}, nil
}

func contentTypeGet(file *multipart.FileHeader) (string, error) {

	mime := file.Header.Get("Content-Type")

	return mime, nil
}

func (fs *FileStorage) Add(file *multipart.FileHeader, typ string, claims *infrastruct.CustomClaims) (*File, error) {
	mime, err := contentTypeGet(file)
	if err != nil {
		return nil, err
	}

	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	uploadInfo, err := fs.client.PutObject(context.Background(), typ, fmt.Sprintf("%s/%d/%s", claims.Role, claims.UserID, file.Filename), f, file.Size, minio.PutObjectOptions{
		ContentType: mime,
	})
	if err != nil {
		return nil, err
	}

	fileObject := &File{
		Url:      file.Filename,
		Length:   uploadInfo.Size,
		Tag:      uploadInfo.Key,
		MimeType: mime,
		Bucket:   uploadInfo.Bucket,
		Object:   file.Filename,
		Role:     claims.Role,
		UserId:   claims.UserID,
	}

	return fileObject, nil
}
