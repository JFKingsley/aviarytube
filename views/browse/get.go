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
	Key string
	HasFull bool
	HasAlternate bool
}

func GetBrowse(c *gin.Context) {
	svc := s3.New(session.New(), &aws.Config{Region: aws.String(configuration.AWSRegion)})

	params := &s3.ListObjectsInput{
		Bucket: aws.String(configuration.Bucket),
	}

	resp, _ := svc.ListObjects(params)

	folders := map[string]*Folder{}

	for _, key := range resp.Contents {
		keyData := *key.Key

		folder := ""
		file := ""

		if strings.HasSuffix(keyData, "/") {
			folder = keyData[:len(keyData) - 1]
		} else {
			splitKey := strings.Split(keyData, "/")
			folder = splitKey[0]
			file = splitKey[1]
		}


		if folders[folder] == nil {
			t, _ := time.Parse("02 January 2006", folder)
			folders[folder] = &Folder{
				Name: folder,
				Key: base64.StdEncoding.EncodeToString([]byte(folder)),
				Date: t,
			}
		}

		if file == "recording.mp4" {
			folders[folder].HasFull = true
		}

		if file == "recording-720.mp4" {
			folders[folder].HasAlternate = true
		}
	}

	folderList := make([]Folder, 0)

	for _, folder := range folders {
		folderList = append(folderList, *folder)
	}

	sort.Slice(folderList, func(i, j int) bool {
		return folderList[i].Date.Before(folderList[j].Date)
	})

	c.HTML(http.StatusOK, "browse/index.html", gin.H{
		"Data": folderList,
	})
}
