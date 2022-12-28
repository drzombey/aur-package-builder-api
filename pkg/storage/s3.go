package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type S3Provider struct {
	client *minio.Client
	cfg    *S3Config
}

type S3Config struct {
	BucketName string
	AccessKey  string
	SecretKey  string
	Endpoint   string
	Region     string
	Ssl        bool
}

func NewS3Provider(cfg *S3Config) Provider {
	ctx := context.Background()

	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Ssl,
	})

	if err != nil {
		log.Fatalln(err)
	}

	err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{Region: cfg.Region})
	if err != nil {
		exists, err := minioClient.BucketExists(ctx, cfg.BucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", cfg.BucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", cfg.BucketName)
	}

	return &S3Provider{
		cfg:    cfg,
		client: minioClient,
	}
}

func (p S3Provider) AddFile(ctx context.Context, fileName string, stream io.Reader) error {
	_, err := createTarFileFromStream(fileName, os.TempDir(), stream)

	if err != nil {
		return err
	}
	_, err = p.client.FPutObject(ctx, p.cfg.BucketName, fileName, os.TempDir()+"/"+fileName, minio.PutObjectOptions{ContentType: "application/zst"})

	if err != nil {
		log.Errorf("Cannot upload following file %s to bucket\n", fileName)
		return err
	}

	return nil
}

func (p S3Provider) GetFile(ctx context.Context, filename string) (*os.File, error) {
	reader, err := p.client.GetObject(ctx, p.cfg.BucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer reader.Close()

	localFile, err := os.Create(filename)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	stat, err := reader.Stat()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
		log.Error(err)
		return nil, err
	}

	return localFile, err
}

func (p S3Provider) DeleteFile(ctx context.Context, filename string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	if err := p.client.RemoveObject(ctx, p.cfg.BucketName, filename, opts); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
