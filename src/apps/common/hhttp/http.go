package hhttp

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
)

func MD5(data string) string {
	m := md5.New()
	m.Write([]byte("$#@$32sfs$%^cs" + data))
	return hex.EncodeToString(m.Sum(nil))
}

func HttpPostBodyWithToken(url, token string, data []byte) ([]byte, error) {
	return HttpPostBodyWithHeader("POST", url, data, map[string]string{"Token": token})
}

func HttpGetBodyWithToken(url, token string, data []byte) ([]byte, error) {
	return HttpPostBodyWithHeader("GET", url, data, map[string]string{"Token": token})
}

func HttpPostBodyWithHeader(method string, url string, data []byte, header map[string]string) ([]byte, error) {
	body := bytes.NewReader(data)
	var resp *http.Response

	httpclient := &http.Client{}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range header {
		request.Header.Set(k, v)
	}
	resp, err = httpclient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
