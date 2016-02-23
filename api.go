package gob2
import (
	"github.com/pquerna/ffjson/ffjson"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"io"
	"errors"
)


type BucketInfo struct {
	AccountId  string `json:"accountId"`
	BucketId   string `json:"bucketId"`
	BucketName string `json:"bucketName"`
	BucketType string `json:"bucketType"`
}

type FileInfo struct {
	FileId string `json:"fileId"`
	FileName string `json:"fileName"`
	Action string `json:"action"`
	Size int `json:"size"`
	UploadTimestamp int `json:"uploadTimestamp"`
}

// GET /b2_authorize_account
// HTTP BASIC "ACCOUNT_ID:APPLICATION_KEY"
type AuthorizeAccountResponse struct {
	AccountId          string `json:"accountId"`
	ApiUrl             string `json:"apiUrl"`
	AuthorizationToken string `json:"authorizationToken"`
	DownloadUrl        string `json:"downloadUrl"`
}

// POST /b2_get_file_info
type GetFileInfoRequest struct {
	FileID string `json:"fileId"`
}
type GetFileInfoResponse struct {
	AccountId     string `json:"accountId"`
	BucketId      string `json:"bucketId"`
	ContentLength int `json:"contentLength"`
	ContentSha1   string `json:"contentSha1"`
	ContentType   string `json:"contentType"`
	FileId        string `json:"fileId"`
	FileInfo      map[string]string `json:"fileInfo"`
	FileName      string `json:"fileName"`
}

// GET /b2_list_buckets
type ListFileNamesRequest struct {
	BucketId string `json:"bucketId"`
	StartFileName string `json:"startFileName"`
	MaxFileCount string `json:"maxFileCount"`
}
type ListFileNamesResponse struct {
	NextFileName string `json:"nextFileName"`
	Files []FileInfo `json:"files"`
}

// POST /b2_list_file_names
type ListBucketsRequest struct {
	AccountId string `json:"accountId"`
}
type ListBucketsResponse struct {
	Buckets []BucketInfo `json:"buckets"`
}


func APIRequest(method string, url string, request interface{}, response interface{}, headers map[string]string) error {
	var reqBodyReader io.Reader = nil
	if request != nil {
		reqBody, err := ffjson.Marshal(request)
		if err != nil {
			return err
		}
		defer ffjson.Pool(reqBody)
		reqBodyReader = bytes.NewReader(reqBody)
	}

	req, err := http.NewRequest(method, url, reqBodyReader)
	if err != nil {
		return err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
	case 401:
		return errors.New(fmt.Sprintf("401 UNAUTHORIZED: %s", body))
	}
	ffjson.Unmarshal(body, response)
	return nil
}
