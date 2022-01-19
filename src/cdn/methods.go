package cdn

import (
	"context"
	"fmt"
	"io"
	"wumpus/src/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func UploadToCdn(file io.Reader, uploadtype string, name string, filetype string) (string, error) {
	ctx := context.Background()
	minioClient, err := minio.New(utils.MINIO_ENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(utils.MINIO_ACCESS_KEY, utils.MINIO_SECRET_KEY, ""),
		Secure: true,
	})

	if err != nil {
		fmt.Println("Failed to create context, got err", err)
		return "", err
	}

	info, uploadErr := minioClient.PutObject(ctx, "wumpcdn", fmt.Sprintf("%s/%s", uploadtype, name), file, -1, minio.PutObjectOptions{ContentType: filetype})
	if uploadErr != nil {
		fmt.Println("Failed to upload file, got err", uploadErr)
		return "", uploadErr
	}

	fmt.Printf("Successfully uploaded %s of size %d\n", name, info.Size)

	return fmt.Sprintf("https://%s/%s/%s", utils.CDN_ENDPOINT, uploadtype, name), nil
}
