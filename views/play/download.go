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

func GetDownload(c *gin.Context) {
	data, err := base64.StdEncoding.DecodeString(c.Param("key"))
	if err != nil {
		log.Fatal("error:", err)
		c.Abort()
		return
	}

	quality := c.Param("quality")

	if quality != "full" && quality != "720p" {
		quality = "full"
	}

	file := "recording.mp4"

	if quality == "720p" {
		file = "recording-720.mp4"
	}

	svc := s3.New(session.New(), &aws.Config{Region: aws.String(configuration.AWSRegion)})

	headObj := s3.HeadObjectInput{
		Bucket: aws.String(configuration.Bucket),
		Key: aws.String(string(data) + "/" + file),
	}

	result, _ := svc.HeadObject(&headObj)

	if *result.ContentType != "video/mp4" {
		c.Status(http.StatusNotFound)
		return
	}

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(configuration.Bucket),
		Key: aws.String(string(data) + "/" + file),
	})

	q := req.HTTPRequest.URL.Query()
	q.Add("response-content-disposition", "attachment; filename=recording.mp4")
	q.Add("response-content-encoding", "video/mp4")
	req.HTTPRequest.URL.RawQuery = q.Encode()

	url, _ := req.Presign(12 * time.Hour)

	c.Redirect(http.StatusTemporaryRedirect, url)
}
