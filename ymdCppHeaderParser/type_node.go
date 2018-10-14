package ymdCppHeaderParser

type Type string

const (
	kPointer    Type = `kPointer`
	kReference  Type = `kReference`
	kLReference Type = `kLReference`
	kLiteral    Type = `kLiteral`
	kTemplate   Type = `kTemplate`
	kFunction   Type = `kFunction`
)

type TypeNode struct {
	IsConst    bool `json:",omitempty"`
	IsVolatile bool `json:",omitempty"`
	IsMutable  bool `json:",omitempty"`
	NodeType   Type `json:",omitempty"`

	// PointerNode
	PointerBase *TypeNode `json:",omitempty"`

	// ReferenceNode
	ReferenceBase *TypeNode `json:",omitempty"`

	// LReferenceNode
	LReferenceBase *TypeNode `json:",omitempty"`

	// TemplateNode
	TemplateBase      *TypeNode   `json:",omitempty"`
	TemplateName      string      `json:",omitempty"`
	TemplateArguments []*TypeNode `json:",omitempty"`

	// LiteralNode
	LiteralName string `json:",omitempty"`

	// FunctionNode
	FunctionReturns   *TypeNode    `json:",omitempty"`
	FunctionArguments [] *Argument `json:",omitempty"`
}

func NewPointerNode(b *TypeNode) *TypeNode {
	return &TypeNode{
		NodeType:    kPointer,
		PointerBase: b,
	}
}

func NewReferenceNode(b *TypeNode) *TypeNode {
	return &TypeNode{
		NodeType:      kReference,
		ReferenceBase: b,
	}
}

func NewLReferenceNode(b *TypeNode) *TypeNode {
	return &TypeNode{
		NodeType:       kLReference,
		LReferenceBase: b,
	}
}

func NewTemplateNode(name string) *TypeNode {
	return &TypeNode{
		NodeType:     kTemplate,
		TemplateName: name,
	}
}

func NewLiteralNode(name string) *TypeNode {
	return &TypeNode{
		NodeType:    kLiteral,
		LiteralName: name,
	}
}

type Argument struct {
	Name string
	Type *TypeNode
}

func NewFunctionNode() *TypeNode {
	return &TypeNode{
		NodeType: kFunction,
	}
}
