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

func GetPlayer(c *gin.Context) {
	data, err := base64.StdEncoding.DecodeString(c.Param("key"))
	if err != nil {
		log.Fatal("error:", err)
		c.Abort()
		return
	}

	svc := s3.New(session.New(), &aws.Config{Region: aws.String(configuration.AWSRegion)})

	headObj := s3.HeadObjectInput{
		Bucket: aws.String(configuration.Bucket),
		Key: aws.String(string(data)),
	}

	result, _ := svc.HeadObject(&headObj)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(configuration.Bucket),
		Key: aws.String(string(data)),
	})

	url, _ := req.Presign(12 * time.Hour)

	c.HTML(http.StatusOK, "play/index.html", gin.H{
		"Name": string(data),
		"Item": result,
		"Link": url,
	})
}
