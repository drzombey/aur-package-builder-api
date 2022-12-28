package storage

import (
	"context"
	"io"
	"os"
)

type Provider interface {
	AddFile(ctx context.Context, filename string, stream io.Reader) error
	GetFile(ctx context.Context, filename string) (*os.File, error)
	DeleteFile(ctx context.Context, filename string) error
}
