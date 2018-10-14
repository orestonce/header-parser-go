package ymdCppHeaderParser

type TokenType string

const (
	kNone       TokenType = `kNone`
	kSymbol     TokenType = `kSymbol`
	kIdentifier TokenType = `kIdentifier`
	kConst      TokenType = `kConst`
)

type ConstType string

const (
	kString  ConstType = `kString`
	kBoolean ConstType = `kBoolean`
	kInt64   ConstType = `kInt64`
	kFloat64 ConstType = `kFloat64`
)

type Token struct {
	MtokenType TokenType `json:",omitempty"`
	MstartPos  int       `json:",omitempty"`
	MstartLine int       `json:",omitempty"`
	Mtoken     string    `json:",omitempty"`

	MconstType ConstType `json:",omitempty"`

	//----union
	MstringConst  string  `json:",omitempty"`
	MboolConst    bool    `json:",omitempty"`
	Mint64Const   int64   `json:",omitempty"`
	Mfloat64Const float64 `json:",omitempty"`
}

type Comment struct {
	text      string
	startLine int
	endLine   int
}
