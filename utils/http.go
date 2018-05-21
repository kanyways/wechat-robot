package utils

import (
	"bytes"
	"net/http"
	"log"
	"io/ioutil"
	"strings"
	"os"
	"mime/multipart"
	"io"
	"fmt"
)

func Get(url string, params map[string]string) ([]byte, error) {
	var r http.Request
	r.ParseForm()
	if len(params) > 0 {
		for key, value := range params {
			r.Form.Add(key, value)
		}
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest(http.MethodGet, url, strings.NewReader(bodystr))
	if err != nil {
		log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		fmt.Println(err)
		return []byte(""), err
	}
	request.Header.Set("Content-Type", "charset=UTF-8")
	request.Header.Set("Connection", "Keep-Alive")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		log.Println("http.Get failed,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("http.Get failed,[err=%s][url=%s]", err, url)
	}
	return b, err
}

func Post(url string, params map[string]string) ([]byte, error) {
	var r http.Request
	r.ParseForm()
	if len(params) > 0 {
		for key, value := range params {
			r.Form.Add(key, value)
		}
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(bodystr))
	if err != nil {
		log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Connection", "Keep-Alive")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		log.Println("http.Post failed,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("http.Post failed,[err=%s][url=%s]", err, url)
	}
	return b, err
}

//body提交二进制数据
func PostBody(url string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	request.Header.Set("Connection", "Keep-Alive")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		log.Println("http.PostBody failed,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("http.PostBody failed,[err=%s][url=%s]", err, url)
	}
	return b, err
}

func PostFile(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request, err
}
