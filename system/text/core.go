package text

var (
	Conv   *ZConv
	Regex  *ZRegex
	File   *ZFile
	String *ZString
)

func init() {
	Conv = NewConv()
	Regex = NewRegex()
	File = NewFile()
	String = NewString()
}
