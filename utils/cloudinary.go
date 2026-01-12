package utils

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func InitCloudinary() error {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return errors.New("cloudinary env variable not set")
	}

	var err error
	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return err
	}

	return nil
}

func UploadToCloudinary(file *multipart.FileHeader) (string, string, error) {
	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	resp, err := cld.Upload.Upload(
		context.Background(),
		f,
		uploader.UploadParams{
			Folder: "adhomes/products",
		},
	)
	if err != nil {
		return "", "", err
	}

	return resp.SecureURL, resp.PublicID, nil
}

func DeleteImageFromCloudinary(publicID string) error {
	if cld == nil {
		return errors.New("cloudinary not initialized")
	}

	fmt.Println("🗑️ Deleting image from Cloudinary:", publicID)

	resp, err := cld.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)

	if err != nil {
		return fmt.Errorf("cloudinary delete error: %v", err)
	}

	fmt.Println("✅ Cloudinary delete result:", resp.Result)
	return nil
}
