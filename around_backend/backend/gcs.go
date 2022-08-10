package backend

import (
	"around/util"
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

var (
	GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

// same as initiate the elasticsearch
func InitGCSBackend(config *util.GCSInfo) {
	//client is like the sessionfactory.openSession() in java, it connects the link with backend
	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	GCSBackend = &GoogleCloudStorageBackend{
		client: client,
		bucket: config.Bucket,
	}
}

// input r io.Reader is the file we upload to GCS, objectName is what this file called
// the return string is the file meidia link that uploaded to GCS, because we need to have this URL record to save into ES
// ES: save the upload record, GCS actually stores the uploaded media
func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
	//create the plain context
	ctx := context.Background()

	//create the empty  object with bucket name and objecname we created (similar to one file that exists in the filesystem)
	object := backend.client.Bucket(backend.bucket).Object(objectName)

	//create the writer context with objected we created
	wc := object.NewWriter(ctx)

	//put actual file r to wc which is the object location, and upload this object to GCS bucket
	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	//ACL=access contro list, let all users have the right to read this file
	//if all have rights to read file, that include the front end, so front end can access this file through URL in GCS
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil
}
