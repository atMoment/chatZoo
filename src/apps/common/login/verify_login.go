package login

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

func SaveLoginToken(cacheUtil db.ICacheUtil, clientPublicKey string, communicationSecret string) error {
	return cacheUtil.Set(redisTableCommunicationKey, combineKey(clientPublicKey, communicationSecret), CommunicationKeyExpiration)
}

func VerifyLoginToken(cacheUtil db.ICacheUtil, clientPublicKey string) (string, error) {
	ret, err := cacheUtil.Get(redisTableCommunicationKey)
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
