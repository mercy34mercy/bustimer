package hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// MD5HashFromData byte[]からMD5によるハッシュ値を算出し文字列にエンコードして返す
func MD5HashFromData(data []byte) (string, error) {
	hash := md5.New()
	if _, err := hash.Write(data); err != nil {
		return "", err
	}
	hashByte := hash.Sum(nil)
	hashString := hex.EncodeToString(hashByte)
	fmt.Println("Hash(MD5): ", hashString)
	return hashString, nil
}
