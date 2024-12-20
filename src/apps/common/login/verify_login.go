package login

// login 服和gate服需要用到里面的函数

import (
	"ChatZoo/common/db"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	redisTableCommunicationKey = "CommunicationKey"
	CommunicationKeyExpiration = 10 * time.Second
)

func GetUserCommunicationKey(userid string) string {
	return fmt.Sprintf("%s:%s", redisTableCommunicationKey, userid)
}

func SaveLoginToken(id string, cacheUtil db.ICacheUtil, clientPublicKey string, communicationSecret string) error {
	return cacheUtil.Set(GetUserCommunicationKey(id), combineKey(clientPublicKey, communicationSecret), CommunicationKeyExpiration)
}

func VerifyLoginToken(id string, cacheUtil db.ICacheUtil, clientPublicKey string) (string, error) {
	if len(clientPublicKey) == 0 {
		return "", errors.New("client public key is empty")
	}
	ret, err := cacheUtil.Get(GetUserCommunicationKey(id))
	if err != nil {
		return "", err
	}
	cpk, ck := splitKey(ret)
	if clientPublicKey != cpk {
		return "", errors.New("client public key not match")
	}
	return ck, nil
}

func combineKey(clientPublicKey string, communicationSecret string) string {
	return fmt.Sprintf("%s-%s", clientPublicKey, communicationSecret)
}

func splitKey(combineKey string) (string, string) {
	ret := strings.Split(combineKey, "-")
	return ret[0], ret[1]
}
