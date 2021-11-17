package tools

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hulasika-fo/zlog/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const logHeader = "HTTP --"

func HttpPostWithHeader(url string, mHead map[string]string) (ret []byte, err error) {
	return httpPost(url, ``, mHead)
}

// HttpGet http get
func HttpGet(url string) (ret []byte, err error) {
	ret, err = httpGet(url)
	if err != nil {
		if strings.Contains(err.Error(), "connection reset") {
			log.Error("请求失败，再次发起请求", err)
			ret, err = httpGet(url)
			if err != nil {
				log.Error("再次发送请求失败", err)
			}
		}
	}
	return
}

func HttpPost(path string) (b []byte, err error) {
	return httpPost(path, ``, nil)
}

func HttpPostWithBodyHeader(path, body string, mHead map[string]string) (b []byte, err error) {
	return httpPost(path, body, mHead)
}

func HttpPostWithBody(path, body string) (b []byte, err error) {
	return httpPost(path, body, nil)
}

func HttpGetWithHeader(url string, mHead map[string]string) (rData []byte, err error) {
	tmId := time.Now().UnixNano()
	log.Debug(tmId, `HttpGetWithHeader -- start`, time.Now().UnixNano(), "url：", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	for k, v := range mHead {
		req.Header.Add(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	rData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(logHeader, "读取反馈失败：", err)
		return
	}
	log.Debug(tmId, `HttpGetWithHeader -- start`, time.Now().UnixNano(), string(rData))
	return
}

// HttpPostFileWithHeader 上传文件
func HttpPostFileWithHeader(path string, mHead map[string]string, file multipart.File, fileName string) (rData []byte, oErr error) {
	oErr = errors.New(`上传文件失败`)
	tmId := time.Now().UnixNano()
	log.Debug(tmId, `HttpPostFileWithHeader -- start`, time.Now().UnixNano(), "url：", path, `header`, mHead)

	boundary := fmt.Sprintf(`----QQQ%v`, time.Now().UnixNano())

	f, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		return
	}

	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)
	_, _ = io.WriteString(buf, fmt.Sprintf("--%v\r\n", boundary))
	_, _ = io.WriteString(buf, fmt.Sprintf(`Content-Disposition: form-data; name="%v"; filename="%v"`, "file", fileName)+"\r\n")
	_, _ = io.WriteString(buf, "Content-Type: video/mp4\r\n\r\n")
	_, _ = io.WriteString(buf, string(f))
	_, _ = io.WriteString(buf, fmt.Sprintf("\r\n--%v--\r\n", boundary))

	req, err := http.NewRequest("POST", path, buf)
	if err != nil {
		log.Error(logHeader, "新建请求失败：", err)
		return
	}
	for k, v := range mHead {
		req.Header.Add(k, v)
	}
	req.Header.Add(`Accept`, `*/*`)
	req.Header.Add(`Content-Type`, fmt.Sprintf(`multipart/form-data; boundary=%v`, boundary))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(logHeader, "发起请求失败：", err)
		return
	}
	defer func() {
		_ = res.Body.Close()
	}()
	rData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(logHeader, "读取反馈失败：", err)
		return
	}
	log.Debug(tmId, `HttpPostFileWithHeader -- end`, time.Now().UnixNano(), string(rData))
	oErr = nil
	return
}

// http get请求
func httpGet(url string) (ret []byte, err error) {
	log.Debug(logHeader, "Get Url：", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	ret, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(logHeader, "读取反馈失败：", err)
		return
	}
	return
}

func httpPost(url, body string, mHead map[string]string) (ret []byte, err error) {
	tmId := time.Now().UnixNano()
	log.Debug(tmId, `HttpPost -- start`, time.Now().UnixNano(), "url：", url, `header`, mHead, "body:", body)
	payload := strings.NewReader(body)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Error(logHeader, "新建请求失败：", err)
		return
	}
	for k, v := range mHead {
		req.Header.Add(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(logHeader, "发起请求失败：", err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	ret, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(logHeader, "读取反馈失败：", err)
		return
	}
	log.Debug(tmId, `HttpPost -- end`, time.Now().UnixNano(), string(ret))
	return
}
