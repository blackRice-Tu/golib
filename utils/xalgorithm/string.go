package xalgorithm

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"github.com/blackRice-Tu/golib/utils/xcommon"

	"github.com/speps/go-hashids/v2"
)

// AlphanumericEncrypt ...
func AlphanumericEncrypt(input string, toLower ...bool) (output string) {
	for _, char := range input {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			output += string(char)
		}
	}
	if len(toLower) == 0 {
		return
	}
	if toLower[0] {
		output = strings.ToLower(output)
	} else {
		output = strings.ToUpper(output)
	}
	return
}

// GetStringMd5 ...
func GetStringMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GetStringSha256 ...
func GetStringSha256(str string, toUpper bool) string {
	h := sha256.New()
	h.Write([]byte(str))
	sign := hex.EncodeToString(h.Sum(nil))
	if toUpper {
		sign = strings.ToUpper(sign)
	}
	return sign
}

// GetRandString ...
func GetRandString(length int) string {
	str := xcommon.Digits + xcommon.AsciiLetters
	bytes := []byte(str)
	var result []byte
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GetRandStringWithPunctuation ...
func GetRandStringWithPunctuation(length int) string {
	str := xcommon.Digits + xcommon.AsciiLetters + xcommon.PunctuationSimple
	bytes := []byte(str)
	var result []byte
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GenerateTimeId() string {
	return time.Now().Format("20060102150405")
}

func HashIdEncodeInt64(id int64, salt string, minLength int) (hashId string) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, _ := hashids.NewWithData(hd)
	hashId, _ = h.EncodeInt64([]int64{id})
	return
}
