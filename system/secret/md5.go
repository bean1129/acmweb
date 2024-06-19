package secret

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"acmweb/system/text"
)

type ZMD5 struct{}

func NewMD5() *ZMD5 {
	return &ZMD5{}
}

func (c *ZMD5) MD5(data string) (string, error) {
	//v, err := c.Encrypt(data)
	//if err != nil {
	//	return "", err
	//}
	//return v, nil
	h1 := md5.New()
	if _, err := h1.Write([]byte(data)); err != nil {
		return "", err
	}
	md5V1 := hex.EncodeToString(h1.Sum(nil))
	return md5V1, nil
}

func (c *ZMD5) Password(data string) (string, error) {
	// 第一次MD5加密
	data1, err := c.Encrypt(data)
	if err != nil {
		return "", err
	}
	// 第二次MD5加密
	data2, err := c.Encrypt(data1)
	if err != nil {
		return "", err
	}
	return data2, nil
}

// Encrypt MD5加密字符串
// encrypts any type of variable using ZMD5 algorithms.
// It uses gconv package to convert <v> to its bytes type.
func (c *ZMD5) Encrypt(data interface{}) (encrypt string, err error) {
	return c.EncryptBytes(text.Conv.Bytes(data))
}

// MustEncrypt encrypts any type of variable using ZMD5 algorithms.
// It uses gconv package to convert <v> to its bytes type.
// It panics if any error occurs.
func (c *ZMD5) MustEncrypt(data interface{}) string {
	result, err := c.Encrypt(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptBytes encrypts <data> using ZMD5 algorithms.
func (c *ZMD5) EncryptBytes(data []byte) (encrypt string, err error) {
	h := md5.New()
	if _, err = h.Write([]byte(data)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncryptBytes MustEncryptBytes encrypts <data> using ZMD5 algorithms.
// It panics if any error occurs.
func (c *ZMD5) MustEncryptBytes(data []byte) string {
	result, err := c.EncryptBytes(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptString EncryptBytes encrypts string <data> using ZMD5 algorithms.
func (c *ZMD5) EncryptString(data string) (encrypt string, err error) {
	return c.EncryptBytes([]byte(data))
}

// MustEncryptString encrypts string <data> using ZMD5 algorithms.
// It panics if any error occurs.
func (c *ZMD5) MustEncryptString(data string) string {
	result, err := c.EncryptString(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptFile encrypts file content of <path> using ZMD5 algorithms.
func (c *ZMD5) EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncryptFile encrypts file content of <path> using ZMD5 algorithms.
// It panics if any error occurs.
func (c *ZMD5) MustEncryptFile(path string) string {
	result, err := c.EncryptFile(path)
	if err != nil {
		panic(err)
	}
	return result
}
