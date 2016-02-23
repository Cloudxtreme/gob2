How to use
=====
```
package main

import (
	"github.com/rikuayanokozy/gob2"
	"fmt"
)

func main() {
	b2, err := gob2.NewB2("ACCOUNT_ID"", "APP_KEY")
	if err != nil {
		panic(err)
	}
	bucket := b2.GetBucketByName("BUCKET_NAME")
	if bucket == nil {
		fmt.Println("Bucket", bucket, "does not exists.")
		return
	}
	files, err := bucket.ListFileNames()
	if err != nil {
		panic(err)
	}
	fmt.Println("Listing b2://" + bucket.BucketName + ":")
	for idx, file := range files {
		fmt.Println("   ", idx, file.FileName)
	}
}
```
