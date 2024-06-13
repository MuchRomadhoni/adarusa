package cloudinary

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloudinary *cloudinary.Cloudinary

func InitCloudinary() {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatalf("failed to initialized Cludinary: %v", err)
	}
	Cloudinary = cld
}