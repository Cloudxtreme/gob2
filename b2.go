package gob2

import (
	"time"
	"encoding/base64"
	"fmt"
	"net/http"
)


type B2 struct {
	ApiUrl                   string
	DownloadUrl              string

	AccountId                string
	AuthorizationToken       string
	AuthorizationTokenExpire int

	Buckets                  []*Bucket
}


func NewB2(accountId string, applicationKey string) (*B2, error) {
	resp := &AuthorizeAccountResponse{}
	headers := make(map[string]string)
	headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(accountId + ":" + applicationKey))
	err := APIRequest("GET", "https://api.backblaze.com/b2api/v1/b2_authorize_account", nil, resp, headers)
	if err != nil {
		return nil, err
	}
	b2 := &B2{
		ApiUrl: resp.ApiUrl + "/b2api/v1/",
		DownloadUrl: resp.DownloadUrl,

		AccountId: resp.AccountId,
		AuthorizationToken: resp.AuthorizationToken,
		AuthorizationTokenExpire: int(time.Now().Unix()) + int(86400),
	}
	b2.ListBuckets()
	return b2, nil
}

func (self *B2) apiRequest(method string, url string, request interface{}, response interface{}, headers map[string]string) error {
	if headers == nil {
		headers = make(map[string]string)
	}
	if request != nil {
		headers["Content-Type"] = "application/json"
	}
	headers["Authorization"] = self.AuthorizationToken
	fmt.Println(method, self.ApiUrl, url, request, response, headers)
	return APIRequest(method, self.ApiUrl + url, request, response, headers)
}

func (self *B2) ListBuckets() ([]*Bucket, error) {
	resp := &ListBucketsResponse{}
	err := self.apiRequest("GET", "b2_list_buckets?accountId=" + self.AccountId, nil, &resp, nil)
	if (err != nil) {
		return nil, err
	}
	self.Buckets = make([]*Bucket, len(resp.Buckets))
	for idx, bucket_info := range resp.Buckets {
		self.Buckets[idx] = &Bucket {
			AccountId: bucket_info.AccountId,
			BucketId: bucket_info.BucketId,
			BucketName: bucket_info.BucketName,
			BucketType: bucket_info.BucketType,
			B2: self,
		}
	}
	return self.Buckets, nil
}

func (self *B2) GetBucketById(id string) *Bucket {
	for _, bucket := range self.Buckets {
		if bucket.BucketId == id {
			return bucket
		}
	}
	return nil
}

func (self *B2) GetBucketByName(name string) *Bucket {
	for _, bucket := range self.Buckets {
		if bucket.BucketName == name {
			return bucket
		}
	}
	return nil
}




type Bucket struct {
	AccountId  string
	BucketId   string
	BucketName string
	BucketType string
	B2 *B2
}
func (self *Bucket) ListFileNames() ([]FileInfo, error) {
	var req *ListFileNamesRequest = nil
	//	req := &ListFileNamesRequest{
	//		BucketId: self.BucketId,
	//	}
	resp := &ListFileNamesResponse{}
	var err error
	if req != nil {
		err = self.B2.apiRequest("POST", "b2_list_file_names", req, resp, nil)
	} else {
		err = self.B2.apiRequest("GET", "b2_list_file_names?bucketId=" + self.BucketId, nil, resp, nil)
	}
	if err != nil {
		return nil, err
	}
	return resp.Files, nil
}

func (self *Bucket) DownloadFileByName(name string) (*http.Response, error) {
	url := self.B2.DownloadUrl + "/file/" + self.BucketName + "/" + name
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", self.B2.AuthorizationToken)
	return http.DefaultClient.Do(req)
}
