package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
)

func Md5Hash(str string) string {
	md := md5.New()
	io.WriteString(md, str)
	return fmt.Sprintf("%x", md.Sum(nil))
}

func GetFileHashName(filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	timeStamp := time.Now().UTC().String()

	hasher := sha1.New()
	_, _ = hasher.Write([]byte(name + timeStamp))

	hashedName := hex.EncodeToString(hasher.Sum(nil)) + ext
	return hashedName
}
