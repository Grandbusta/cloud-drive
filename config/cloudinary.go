package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCld() *cloudinary.Cloudinary {
	cld, err := cloudinary.New()
	if err != nil {
		log.Fatal(err)
	}
	return cld
}
