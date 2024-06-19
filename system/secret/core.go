package secret

var (
	MD5 *ZMD5
	DES *ZDES
)

func init() {
	MD5 = NewMD5()
	DES = NewZDES()
}
