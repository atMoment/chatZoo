package main

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	Code_Success = 1
	Code_Failed  = 2
)

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
