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
	// Initialize Cloudinary
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	// Delete the image
	_, err = cld.Upload.Destroy(ctx(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from Cloudinary: %v", err)
	}

	return nil
}

// ctx is a helper to create a context for Cloudinary API calls
func ctx() (c context.Context) {
	c = context.Background()
	return
}
