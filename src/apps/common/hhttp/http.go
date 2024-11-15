package hhttp

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

func ParseJsonReq(r *http.Request, req interface{}) error {
	//tk := r.Header.Get("token")
	//if len(tk) > 0 && tk != token {
	//	return fmt.Errorf("token not match")
	//}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, req)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonRsp(w http.ResponseWriter, rsp interface{}) error {
	b, err := json.Marshal(rsp)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}
