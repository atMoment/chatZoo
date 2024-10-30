package encrypt

import (
	"crypto/rc4"
	"errors"
)

// 消息使用通信密钥加密解密

type ICoder interface {
	Encryption(data []byte, afterEncryptData []byte) error
	Decryption(data []byte, afterDecryptData []byte) error
}

type _Coder struct {
	key []byte // 通信密钥
}

func NewCoder(key []byte) ICoder {
	if len(key) == 0 || len(key) > 256 {
		panic("key illegal")
	}
	return &_Coder{
		key: key,
	}
}

// Encryption 对消息使用通信密钥 加密
func (c *_Coder) Encryption(data []byte, afterEncryptData []byte) error {
	if cap(afterEncryptData) < len(data) {
		return errors.New("no enough space")
	}
	rc, err := rc4.NewCipher(c.key)
	if err != nil {
		return err
	}
	// data 和 afterEncryptData 底层数组范围要么完全一样, 要么完全不一样
	rc.XORKeyStream(data, afterEncryptData)
	return nil
}

// Decryption 对消息使用通信密钥解密
func (c *_Coder) Decryption(data []byte, afterDecryptData []byte) error {
	// 再走一遍加密流程就是解密
	return c.Encryption(data, afterDecryptData)
}
