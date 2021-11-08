package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/15125505/zlog/log"
	"os/exec"
	"runtime"
	"strings"
)

func StringBetween(src, first, second string) (ret string, err error) {
	pos := strings.Index(src, first)
	if -1 == pos {
		err = errors.New("first字段不存在")
		return
	}
	src = src[pos+len(first):]
	if len(second) == 0 {
		ret = src
		return
	}
	pos = strings.Index(src, second)
	if -1 == pos {
		err = errors.New("second字段不存在")
		return
	}
	ret = src[:pos]
	return
}

func Concat(s []string, split string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(s), "[]"), " ", split, -1)
}

func AesDecrypt(decodeStr, sessionKey, ivStr string) ([]byte, error) {

	//先解密base64
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}

	iv, err := base64.StdEncoding.DecodeString(ivStr)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 去重复
func IntSliceDereplication(s []int) []int {
	m := make(map[int]bool)

	for i := 0; i < len(s); i++ {
		if _, ok := m[s[i]]; ok {
			s = append(s[:i], s[i+1:]...)
			i--
			continue
		}
		m[s[i]] = true
	}
	return s
}

func Int64SliceDereplication(s []int64) []int64 {
	m := make(map[int64]bool)

	for i := 0; i < len(s); i++ {
		if _, ok := m[s[i]]; ok {
			s = append(s[:i], s[i+1:]...)
			i--
			continue
		}
		m[s[i]] = true
	}
	return s
}

func StringSliceDereplication(s []string) []string {
	m := make(map[string]bool)

	for i := 0; i < len(s); i++ {
		if _, ok := m[s[i]]; ok {
			s = append(s[:i], s[i+1:]...)
			i--
			continue
		}
		m[s[i]] = true
	}
	return s
}

func GetSep() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

func DoCommand(cmd string, stdout, stderr *bytes.Buffer) (err error) {
	c := "cmd"
	a := "/c"
	if strings.ToLower(runtime.GOOS) == "linux" {
		c = "bash"
		a = "-c"
	}
	log.Info(`exec command`, c, a, cmd)
	p := exec.Command(c, a, cmd)
	if stdout != nil {
		p.Stdout = stdout
	}
	if stderr != nil {
		p.Stderr = stderr
	}
	err = p.Run()
	return
}
