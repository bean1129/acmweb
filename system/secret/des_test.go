package secret

import (
	"bytes"
	"fmt"
	"testing"
)

type destTest struct {
	cleartxt  string
	ciphertxt []byte
}

var golden = []destTest{
	{"aabbccddff",
		[]byte{36, 122, 85, 107, 3, 196, 56, 234, 178, 93, 252, 40, 18, 222, 146, 65},
	},
}

var key = "12345678"

func TestDesEncryopt(t *testing.T) {
	DES := NewZDES()
	for _, obj := range golden {
		result, _ := DES.DesEncryopt([]byte(obj.cleartxt), []byte(key))
		if bytes.Compare(result, obj.ciphertxt) == 0 {
			fmt.Println("Encryopt success")
		} else {
			fmt.Println(result)
		}
	}
}
func TestDesDecrypt(t *testing.T) {
	DES := NewZDES()
	src := []byte{36, 122, 85, 107, 3, 196, 56, 234, 178, 93, 252, 40, 18, 222, 146, 65}
	key := "12345678"
	result, _ := DES.DesDecrypt(src, []byte(key))
	fmt.Println(string(result))
}
