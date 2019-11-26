package timetable

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"cloud.google.com/go/storage"
)

// fetchFromCloudStorage CloudStorageからfileNameに該当するファイルのbyteデータを取得する
func fetchFromCloudStorage(fileName string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return []byte{}, err
	}
	bkt := client.Bucket(bucket)
	obj := bkt.Object(fileName)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return []byte{}, err
	}
	defer r.Close()
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// saveToCloudStorage dataのバイト列をfileNameで示されるファイルに保存する
func saveToCloudStorage(data []byte, fileName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	bkt := client.Bucket(bucket)
	obj := bkt.Object(fileName)
	w := obj.NewWriter(ctx)
	if _, err := w.Write(data); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func saveCache(data interface{}, fileName string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = saveToCloudStorage(jsonData, fileName)
	if err != nil {
		return err
	}
	return nil
}
