package s3

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	sdkaws "github.com/blend/go-sdk/aws"
)

// File is info for a file upload.
type File struct {
	FilePath    string
	Bucket      string
	Key         string
	ContentType string
}

// UploadFile uploads a file to s3.
func UploadFile(ctx context.Context, cfg sdkaws.Config, fileInfo File) error {
	// Open the file for use
	file, err := os.Open(fileInfo.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return err
	}
	size := int64(stats.Size())

	session := sdkaws.MustNewSession(cfg)
	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(fileInfo.Bucket),
		Key:                  aws.String(fileInfo.Key),
		ACL:                  aws.String("private"),
		Body:                 file,
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(fileInfo.ContentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
