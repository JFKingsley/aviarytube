package play

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jfkingsley/aviarytube/configuration"
	"log"
	"net/http"
	"time"
)

type PlayerData struct {
	FullURL string
	AlternateURL string
	DownloadURL string
}

func GetPlayer(c *gin.Context) {
	data, err := base64.StdEncoding.DecodeString(c.Param("key"))
	if err != nil {
		log.Fatal("error:", err)
		c.Abort()
		return
	}

	svc := s3.New(session.New(), &aws.Config{Region: aws.String(configuration.AWSRegion)})

	c.HTML(http.StatusOK, "play/index.html", gin.H{
		"Name": string(data),
		"Data": &PlayerData{
			FullURL: getSignedURL(svc, string(data), "recording.mp4"),
			AlternateURL: getSignedURL(svc, string(data), "recording-720.mp4"),
			DownloadURL: "/play/" + c.Param("key") + "/download/full",
		},
	})
}

func getSignedURL(svc *s3.S3, folder string, file string) string {
	headObj := s3.HeadObjectInput{
		Bucket: aws.String(configuration.Bucket),
		Key: aws.String(folder + "/" + file),
	}

	result, _ := svc.HeadObject(&headObj)

	if *result.ContentType != "video/mp4" {
		return ""
	}

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(configuration.Bucket),
		Key: aws.String(folder + "/" + file),
	})

	url, _ := req.Presign(12 * time.Hour)

	return url
}
