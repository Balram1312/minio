package utils

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type App struct {
	MinioClient *minio.Client
	Context     context.Context
}
