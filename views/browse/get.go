package browse

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jfkingsley/aviarytube/configuration"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Item struct {
	Name string
	Key string
}

type Folder struct {
	Name string
	Date time.Time
	Items []Item
}

func GetBrowse(c *gin.Context) {
	svc := s3.New(session.New(), &aws.Config{Region: aws.String(configuration.AWSRegion)})

	params := &s3.ListObjectsInput{
		Bucket: aws.String(configuration.Bucket),
	}

	resp, _ := svc.ListObjects(params)

	data := map[string][]Item{}

	for _, key := range resp.Contents {
		keyData := *key.Key
		if strings.Contains(keyData, "/") {
			splitKey := strings.Split(keyData, "/")
			if data[splitKey[0]] == nil {
				data[splitKey[0]] = []Item{{splitKey[1], base64.StdEncoding.EncodeToString([]byte(keyData))}}
			} else {
				data[splitKey[0]] = append(data[splitKey[0]], Item{splitKey[1], base64.StdEncoding.EncodeToString([]byte(keyData))})
			}
		}
	}

	folders := make([]Folder, 0)

	for key, value := range data {
		t, _ := time.Parse("02 January 2006", key)
		folders = append(folders, Folder{
			Name: key,
			Date: t,
			Items: value,
		})
	}

	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Date.Before(folders[j].Date)
	})

	c.HTML(http.StatusOK, "browse/index.html", gin.H{
		"Data": folders,
	})
}
