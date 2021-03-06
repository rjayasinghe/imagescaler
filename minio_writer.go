package main

import (
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v6"
	uuid "github.com/nu7hatch/gouuid"
)

//TODO: create an interface and have several implementations e.g. test mocks

func writeImageToObjectStorage(scaledReader io.Reader, length int, imageType string, targetScale string, config imageScalerConfig) (ImageUpdate, error) {
	var imageUpdate ImageUpdate
	var minioOpts minio.PutObjectOptions
	minioOpts.ContentType = fmt.Sprintf("image/%s", imageType)
	imageUUID, err := uuid.NewV4()
	if err != nil {
		log.Printf("error while creating image UUID: %v", err)
		return imageUpdate, err
	}
	minioClient, err := minio.New(config.minioURL, config.minioAccessKey, config.minioSecret, config.minioSecure)
	if err != nil {
		log.Printf("error while creating min.io client: %v", err)
		return imageUpdate, err
	}
	_, err = minioClient.PutObject(config.minioBucketName, imageUUID.String(), scaledReader, -1, minioOpts)
	if err != nil {
		log.Printf("error while writing image to min.io: %v", err)
		return imageUpdate, err
	}

	imageUpdate.ImageUUID = imageUUID.String()
	imageUpdate.URL = fmt.Sprintf("%s/%s/%s", config.minioExternalURL, config.minioBucketName, imageUUID)
	if err != nil {
		return imageUpdate, err
	}
	imageUpdate.ImageScale = targetScale

	return imageUpdate, nil
}
