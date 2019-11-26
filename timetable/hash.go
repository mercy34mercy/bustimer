package timetable

import (
	"github.com/shun-shun123/bus-timer/hash"
)

// md5HashFromGCSFile Google Cloud Storageに保存されているファイルからハッシュ値を算出する
func md5HashFromGCSFile(fileName string) (string, error) {
	data, err := fetchFromCloudStorage(fileName)
	if err != nil {
		return "", err
	}
	hashString, err := hash.MD5HashFromData(data)
	if err != nil {
		return "", err
	}
	return hashString, nil
}
