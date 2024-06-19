package secret

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"strconv"
)

//临时
type DESCError int

func (d DESCError) Error() string {
	return "need a multiple of the blocksize :" + strconv.Itoa(int(d))
}

type ZDES struct{}

func NewZDES() *ZDES {
	return &ZDES{}
}

func (c *ZDES) padding(src []byte, blocksize int) []byte {
	n := len(src)
	padnum := blocksize - n%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	dst := append(src, pad...)
	return dst
}

func (c *ZDES) unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	dst := src[:n-unpadnum]
	return dst
}

func (c *ZDES) DesEncryopt(src []byte, key []byte) ([]byte, error) {
	key = key[:8]
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = c.padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return src, nil
}

func (c *ZDES) DesDecrypt(src []byte, key []byte) ([]byte, error) {
	key = key[:8]
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = c.unpadding(src)
	return src, nil
}

//func (c *ZDES) DesECBEncrypt(data, key []byte) ([]byte, error) {
//	//NewCipher创建一个新的加密块
//	key = key[:8]
//	block, err := des.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//
//	bs := block.BlockSize()
//	data = c.Pkcs5Padding(data, bs)
//	if len(data)%bs != 0 {
//		return nil, DESCError(len(data))
//	}
//
//	out := make([]byte, len(data))
//	dst := out
//	for len(data) > 0 {
//		//Encrypt加密第一个块，将其结果保存到dst
//		block.Encrypt(dst, data[:bs])
//		data = data[bs:]
//		dst = dst[bs:]
//	}
//	return out, nil
//}
//
//func (c *ZDES) Pkcs5Padding(ciphertext []byte, blockSize int) []byte {
//	padding := blockSize - len(ciphertext)%blockSize
//	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(ciphertext, padtext...)
//}
