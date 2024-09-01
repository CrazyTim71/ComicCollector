package operations

import (
	"bytes"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func GetImageById(bucket *gridfs.Bucket, id primitive.ObjectID) (bytes.Buffer, error) {
	var buf bytes.Buffer
	_, err := bucket.DownloadToStream(id, &buf)

	return buf, err
}

func GetImageByFilename(bucket *gridfs.Bucket, filename string) (bytes.Buffer, error) {
	var buf bytes.Buffer
	_, err := bucket.DownloadToStreamByName(filename, &buf)

	return buf, err
}

func DeleteImage(bucket *gridfs.Bucket, id primitive.ObjectID) error {
	return bucket.Delete(id)
}
