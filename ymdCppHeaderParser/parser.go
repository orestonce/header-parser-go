package ymdCppHeaderParser

import (
	"log"
)

type ScopeType string

const (
	kGlobal    ScopeType = `kGlobal`
	kNamespace ScopeType = `kNamespace`
	kClass     ScopeType = `kClass`
)

type AccessControlType string

const (
	kPublic    AccessControlType = `kPublic`
	kPrivate   AccessControlType = `kPrivate`
	kProtected AccessControlType = `kProtected`
)

type Scope struct {
	scopeType                ScopeType
	name                     string
	currentAccessControlType AccessControlType
}

type Parser struct {
	Tokenizer
	scopes      [64]Scope
	topScopeIdx int

	debug bool
}

func NewParser(input []byte) *Parser {
	this := &Parser{}

	// Pass the input to the tokenizer
	this.Tokenizer = *NewTokenizer(input, 1)
	// Reset scope
	topScope := this.scopes[this.topScopeIdx]
	topScope.name = ``
	topScope.scopeType = kGlobal
	topScope.currentAccessControlType = kPublic
	return this
}

func (this *Parser) ParseAll() {
	// Parse all statements in the file
	for this.ParseStatement() {
	}
}

func (this *Parser) ParseStatement() bool {
	var token Token
	if !this.GetToken(&token, false, false) {
		return false
	}
	if !this.ParseDeclaration(&token) {
		return false
	}

	return true
}

func (this *Parser) ParseDeclaration(token *Token) bool {
	const funcId = `8w3c6jsa `

	this.debugPrintf(funcId, "token %v", marshalJson(token))

	switch token.Mtoken {
	case `#`:
		this.UngetToken(token)
		return this.ParseDirective()
	case `namespace`:
		this.UngetToken(token)
		return this.ParseNamespace()
	case `;`:
		return true
	case `enum`:
		this.UngetToken(token)
		this.ParseEnum()
		return true
	case `class`, `struct`:
		this.UngetToken(token)
		return this.ParseClass()
	case `template`:
		return this.SkipDeclaration(token)
	}
	if this.ParseAccessControl(token, &this.scopes[this.topScopeIdx].currentAccessControlType) {
		this.RequireSymbol(`:`)
		this.debugPrintf(funcId, "is access control")
		return true
	}

	this.UngetToken(token)
	// function or property ?
	// Process method specifiers in any particular order
	isVirtual := false   // method
	isInline := false    // method
	isConstExpr := false // method
	isStatic := false    // method & property

	isMutable := false // property
	for {
		if !isVirtual && this.MatchIdentifier(`virtual`) {
			isVirtual = true
		} else if !isInline && this.MatchIdentifier(`inline`) {
			isInline = true
		} else if !isConstExpr && this.MatchIdentifier(`constexpr`) {
			isConstExpr = true
		} else if !isStatic && this.MatchIdentifier(`static`) {
			isStatic = true
		} else if !isMutable && this.MatchIdentifier(`mutable`) {
			isMutable = true
		} else {
			break
		}
	}

	// Parse the type
	if !this.ParseType() {
		return false
	}

	var nameToken Token
	if !this.GetIdentifier(&nameToken) {
		this.panicf(funcId, `Expected a property or method name`)
	}
	this.debugPrintf(funcId, "nameToken %v", marshalJson(nameToken))

	var next Token
	if !this.GetToken(&next, false, false) {
		panic(`ekuqkxjtz9`)
	}
	switch next.Mtoken {
	case `;`: // is property
		this.debugPrintf(funcId, "token is property")
		return true
	case `(`: // is method
		this.debugPrintf(funcId, "token is method")
		this.UngetToken(token)
		return this.ParseFunction()
	}
	this.debugPrintf(funcId, "skip unknown token")
	return this.SkipDeclaration(token);
}

func (this *Parser) ParseDirective() bool {
	const funcId = `f4haccj6 `
	var token Token

	this.RequireSymbol(`#`)

	// Check the compiler directive
	if !this.GetIdentifier(&token) {
		this.panicf(funcId, `Missing compiler directive after #`)
	}
	this.debugPrintf(funcId, "token %v", marshalJson(token))

	multiLineEnabled := false
	switch token.Mtoken {
	case `define`:
		multiLineEnabled = true
	case `include`:
		var includeToken Token
		this.GetToken(&includeToken, true, false)
		this.debugPrintf(funcId, "includeToken %v", marshalJson(includeToken))
	}

	// Skip past the end of the token
	var lastChar byte = '\n'
	for {
		// Skip to the end of the line
		var c byte
		for {
			if this.is_eof() {
				break
			}
			c = this.GetChar()
			if c == '\n' {
				break
			}
			lastChar = c
		}

		if multiLineEnabled && lastChar == '\\' {
			continue
		} else {
			break
		}
	}

	return true
}

func (this *Parser) SkipDeclaration(token *Token) bool {
	scopeDepth := 0
	for this.GetToken(token, false, false) {
		if token.Mtoken == `;` && scopeDepth == 0 {
			break
		}

		if token.Mtoken == `{` {
			scopeDepth++
		}

		if token.Mtoken == `}` {
			scopeDepth--
			if scopeDepth == 0 {
				break
			}
		}
	}

	return true
}

func (this *Parser) ParseEnum() {
	const funcId = `d3gpz066 `
	if !this.MatchIdentifier(`enum`) {
		this.panicf(funcId, `require "enum" identifier`)
	}
	// C++1x enum class type?
	isEnumClass := this.MatchIdentifier(`class`)

	this.debugPrintf(funcId, "isEnumClass %v", isEnumClass)

	// Parse enum name
	var enumToken Token
	if !this.GetIdentifier(&enumToken) {
		this.panicf(funcId, "Missing enum name")
	}

	this.debugPrintf(funcId, "enum %v", marshalJson(enumToken))

	// Parse C++1x enum base
	if isEnumClass && this.MatchSymbol(`:`) {
		var baseToken Token
		if !this.GetIdentifier(&baseToken) {
			this.panicf(funcId, "Missing enum type specifier")
		}
		// Validate base token
	}

	// Require opening brace
	this.RequireSymbol(`{`)

	// Parse all the values
	var token Token
	for this.GetIdentifier(&token) {
		this.debugPrintf(funcId, "enum object %v", marshalJson(token))
		// Parse constant
		if this.MatchSymbol(`=`) {
			// Just parse the value, not doing anything with it atm
			value := ``
			for this.GetToken(&token, false, false) &&
				(token.MtokenType != kSymbol || (token.Mtoken != `,` && token.Mtoken != `}`)) {
				value += token.Mtoken
			}
			this.debugPrintf(funcId, "value %v", marshalJson(value))
			this.UngetToken(&token)
		}
		// Next value?
		if !this.MatchSymbol(`,`) {
			break
		}
	}

	this.RequireSymbol(`}`)
	this.RequireSymbol(`;`)
}

func (this *Parser) ParseMacroMeta() bool {
	this.RequireSymbol(`(`)
	if !this.ParseMetaSequence() {
		return false
	}
	// Possible ;
	this.MatchSymbol(`;`)

	return true
}

func (this *Parser) ParseMetaSequence() bool {
	const funcId = `hky9quq0 `
	if !this.MatchSymbol(`)`) {
		for {
			// Parse key value
			var keyToken Token
			if !this.GetIdentifier(&keyToken) {
				this.panicf(funcId, "Expected identifier in meta sequence")
			}

			// Simple value?
			if this.MatchSymbol(`=`) {
				var token Token
				if !this.GetToken(&token, false, false) {
					this.panicf(funcId, `Expected token`)
				}
			} else if (this.MatchSymbol(`(`)) { // Compound value
				if !this.ParseMetaSequence() {
					return false
				}
				// No value
			} else {
				// null
			}

			if !this.MatchSymbol(`,`) {
				break
			}
		}

		this.MatchSymbol(`)`)
	}

	return true
}

func (this *Parser) PushScope(name string, scopeType ScopeType, accessControlType AccessControlType) {
	const funcId = `njxt77ngz9 `
	if this.topScopeIdx >= len(this.scopes)-1 {
		this.panicf(funcId, `Max scope depth`)
	}

	this.topScopeIdx++
	topScope := this.scopes[this.topScopeIdx]
	topScope.scopeType = scopeType
	topScope.name = name
	topScope.currentAccessControlType = accessControlType
}

func (this *Parser) PopScope() {
	const funcId = `v83qqwx728 `
	if this.topScopeIdx == 0 {
		this.panicf(funcId, `Scope error`)
	}

	this.topScopeIdx--
}

func (this *Parser) ParseNamespace() bool {
	const funcId = `l4u2kamr `
	var token Token

	if !this.GetIdentifier(&token) || token.Mtoken != `namespace` {
		this.panicf(funcId, `Missing "namespace" identifier`)
	}

	if !this.GetIdentifier(&token) {
		this.panicf(funcId, "Missing namespace name")
	}

	this.RequireSymbol(`{`)

	this.PushScope(token.Mtoken, kNamespace, kPublic)

	for !this.MatchSymbol(`}`) {
		if !this.ParseStatement() {
			return false
		}
	}

	this.PopScope()
	return true
}

func (this *Parser) ParseAccessControl(token *Token, accessControlType *AccessControlType) bool {
	switch token.Mtoken {
	case `public`:
		*accessControlType = kPublic
	case `protected`:
		*accessControlType = kProtected
	case `private`:
		*accessControlType = kPrivate
	default:
		return false
	}
	return true
}

func (this *Parser) ParseClass() bool {
	const funcId = `z0dnwwg6 `

	var startAccessControlType = kPrivate

	if this.MatchIdentifier(`class`) {
		startAccessControlType = kPrivate
	} else if this.MatchIdentifier(`struct`) {
		startAccessControlType = kPublic
	} else {
		this.panicf(funcId, `Missing "class" or "struct"`)
	}
	// Get the class name
	var classNameToken Token
	if !this.GetIdentifier(&classNameToken) {
		this.panicf(funcId, `Missing class name `)
	}
	this.debugPrintf(funcId, "class begin %v", marshalJson(classNameToken))

	if this.MatchSymbol(`;`) { // forward declaration
		this.debugPrintf(funcId, `forward declaration.`)
		return true
	}

	// Match base types
	if this.MatchSymbol(`:`) {
		for {
			var accessOrName Token
			if !this.GetIdentifier(&accessOrName) {
				this.panicf(funcId, `Missing class or access control specifier`)
			}

			// Parse the access control specifier
			accessControlType := startAccessControlType
			if !this.ParseAccessControl(&accessOrName, &accessControlType) {
				this.UngetToken(&accessOrName)
			}
			var baseClassName Token
			if !this.GetIdentifier(&baseClassName) {
				this.panicf(funcId, `Missing base class name`)
			}

			this.debugPrintf(funcId, "base class %v", marshalJson(baseClassName))

			if !this.MatchSymbol(`,`) {
				break
			}
		}
	}

	this.RequireSymbol(`{`)

	this.PushScope(classNameToken.Mtoken, kClass, startAccessControlType)

	for !this.MatchSymbol(`}`) {
		if !this.ParseStatement() {
			return false
		}
	}
	this.PopScope()

	this.RequireSymbol(`;`)
	this.debugPrintf(funcId, "class end %v", marshalJson(classNameToken))

	return true
}

func (this *Parser) ParseProperty(token *Token) bool {
	const funcId = `ajqd8r4p4b `
	if !this.ParseMacroMeta() {
		return false
	}
	// Process method specifiers in any particular order
	isMutable := false
	isStatic := false
	for {
		if !isMutable && this.MatchIdentifier(`mutable`) {
			isMutable = true
		} else if !isStatic && this.MatchIdentifier(`static`) {
			isStatic = true
		} else {
			break
		}
	}
	// Parse the type
	if !this.ParseType() {
		return false
	}

	var nameToken Token
	if !this.GetIdentifier(&nameToken) {
		this.panicf(funcId, `Expected a property name`)
	}

	// Skip until the end of the definition
	var t Token
	for this.GetToken(&t, false, false) {
		if t.Mtoken == `;` {
			break
		}
	}

	return true
}

func (this *Parser) ParseFunction() bool {
	const funcId = `tom77xqc `
	funcNode := NewFunctionNode()
	// Process method specifiers in any particular order
	isVirtual := false
	isInline := false
	isConstExpr := false
	isStatic := false
	for {
		if !isVirtual && this.MatchIdentifier(`virtual`) {
			isVirtual = true
		} else if !isInline && this.MatchIdentifier(`inline`) {
			isInline = true
		} else if !isConstExpr && this.MatchIdentifier(`constexpr`) {
			isConstExpr = true
		} else if !isStatic && this.MatchIdentifier(`static`) {
			isStatic = true
		} else {
			break
		}
	}

	this.debugPrintf(funcId, "front property isVirtual=[%v] isInline=[%v] isConstExpr=[%v] isStatic=[%v]",
		isVirtual, isInline, isConstExpr, isStatic)

	// Parse the return type
	funcNode.FunctionReturns = this.ParseTypeNode()
	this.debugPrintf(funcId, "retType %v", marshalJson(funcNode.FunctionReturns))
	if funcNode.FunctionReturns == nil {
		return false
	}
	// Parse the name of the method
	var nameToken Token
	if !this.GetIdentifier(&nameToken) {
		this.panicf(funcId, `Expected method name`)
	}
	this.debugPrintf(funcId, "function name %v", marshalJson(nameToken))
	this.MatchSymbol("(")
	// Is there an argument list in the first place or is it closed right away?
	if !this.MatchSymbol(`)`) {
		// Walk over all arguments
		for i := 0; ; i++ {
			// Get the type of the argument
			var argTypeNode = this.ParseTypeNode()
			this.debugPrintf(funcId, "argTypeNode %v", marshalJson(argTypeNode))
			if argTypeNode == nil {
				return false
			}
			// Optional argument name
			ok := this.GetIdentifier(&nameToken)
			if ok {
				this.debugPrintf(funcId, "argument name %v", marshalJson(nameToken))
			} else {
				this.debugPrintf(funcId, "argument no name")
			}
			funcNode.FunctionArguments = append(funcNode.FunctionArguments, &Argument{
				Name: nameToken.Mtoken,
				Type: argTypeNode,
			})
			// Parse default value
			if this.MatchSymbol(`=`) {
				defaultValue := ``
				var token Token
				this.GetToken(&token, false, false)
				if token.MtokenType == kConst {
					this.debugPrintf(funcId, "argument default value const %v", marshalJson(token))
				} else {
					for {
						if token.Mtoken == `,` || token.Mtoken == `)` {
							this.UngetToken(&token)
							break
						}
						defaultValue += token.Mtoken
						if !this.GetToken(&token, false, false) {
							break
						}
					}
					this.debugPrintf(funcId, "argument default value express %v", defaultValue)
				}
			} else {
				this.debugPrintf(funcId, "argument have not default value")
			}
			if !this.MatchSymbol(`,`) { // Only in case another is expected
				break
			}
		}
		this.RequireSymbol(`)`)
	}

	// Optionally parse constness
	ok := this.MatchIdentifier(`const`)
	this.debugPrintf(funcId, "function is const %v", ok)
	// Pure?
	if this.MatchSymbol(`=`) {
		var token Token
		if !this.GetToken(&token, false, false) || token.Mtoken != `0` {
			this.panicf(funcId, `Expected nothing else than null`) //
		}
		this.debugPrintf(funcId, `pure func `, marshalJson(token))
	}
	// Skip either the ; or the body of the function
	var skipToken Token
	if !this.SkipDeclaration(&skipToken) {
		return false
	}

	return true
}

func (this *Parser) ParseType() bool {
	node := this.ParseTypeNode()
	if node == nil {
		return false
	}
	return true
}

func (this *Parser) ParseTypeNode() *TypeNode {
	const funcId = `f98vawvz `
	var node *TypeNode
	var token Token

	var (
		isConst    = false
		isVolatile = false
		isMutable  = false
	)
	for {
		if !isConst && this.MatchIdentifier(`const`) {
			isConst = true
		} else if !isVolatile && this.MatchIdentifier(`volatile`) {
			isVolatile = true
		} else if !isMutable && this.MatchIdentifier(`mutable`) {
			isMutable = true
		} else {
			break
		}
	}

	// Parse a literal value
	declarator := ``

	declarator = this.ParseTypeNodeDeclarator()

	// Postfix const specifier
	if this.MatchIdentifier(`const`) {
		isConst = true
	}

	// Template?
	if this.MatchSymbol(`<`) {
		templateNode := NewTemplateNode(declarator)
		for {
			node := this.ParseTypeNode()
			if node == nil {
				return nil
			}
			templateNode.TemplateArguments = append(templateNode.TemplateArguments, node)

			if this.MatchSymbol(`,`) {
				continue
			} else {
				break
			}
		}

		if !this.MatchSymbol(`>`) {
			this.panicf(funcId, `Expected closing > `)
		}
		node = templateNode 
	} else {
		node = NewLiteralNode(declarator)
	}

	// Store gathered stuff
	node.IsConst = isConst

	// Check reference or pointer types
	for this.GetToken(&token, false, false) {
		if token.Mtoken == `&` {
			node = NewReferenceNode(node)
		} else if token.Mtoken == `&&` {
			node = NewLReferenceNode(node)
		} else if token.Mtoken == `*` {
			node = NewPointerNode(node)
		} else {
			this.UngetToken(&token)
			break
		}

		if this.MatchIdentifier(`const`) {
			node.IsConst = true
		}
	}

	// Function pointer?
	if this.MatchSymbol(`(`) {
		// Parse void(*)(args, ...)
		//            ^
		//            |
		if this.MatchSymbol(`*`) {
			var token Token
			this.GetToken(&token, false, false)
			if token.Mtoken != `)` || token.MtokenType != kIdentifier || !this.MatchSymbol(`)`) {
				this.panicf(funcId, `Expected ")"`)
			}
		}

		// Parse arguments
		funcNode := NewFunctionNode()
		funcNode.FunctionReturns = node
		if !this.MatchSymbol(`)`) {
			for {
				argument := Argument{}
				argument.Type = this.ParseTypeNode()
				if argument.Type == nil {
					return nil
				}
				// Get , or name identifier
				if !this.GetToken(&token, false, false) {
					this.panicf(funcId, `Unexpected end of file`)
				}

				// Parse optional name
				if token.MtokenType == kIdentifier {
					argument.Name = token.Mtoken
				} else {
					this.UngetToken(&token)
				}

				funcNode.FunctionArguments = append(funcNode.FunctionArguments, &argument)
				if this.MatchSymbol(`,`) {
					continue
				} else {
					break
				}
			}
			if !this.MatchSymbol(`)`) {
				this.panicf(funcId, `Missing ")"`)
			}
		}

		node = funcNode
	}

	// This stuff refers to the top node
	node.IsVolatile = isVolatile
	node.IsMutable = isMutable

	return node
}

func (this *Parser) ParseTypeNodeDeclarator() string {
	const funcId = `grns8napbd `
	// Skip optional forward declaration specifier
	this.MatchIdentifier(`class`)
	this.MatchIdentifier(`struct`)
	this.MatchIdentifier(`typename`)

	// Parse a type name
	declarator := ``
	var token Token
	first := true
	for { // ::namespace1::namespace2::ClassType
		// Parse the declarator
		if this.MatchSymbol(`::`) {
			declarator += "::"
		} else if !first {
			break
		}

		// Mark that this is not the first time in this loop
		first = false
		// Match an identifier
		if !this.GetIdentifier(&token) {
			this.panicf(funcId, `Expected identifier`)
		}

		declarator += token.Mtoken
	}

	return declarator
}

func (this *Parser) debugPrintf(funcId string, sfmt string, a ... interface{}) {
	if !this.debug {
		return
	}
	log.Printf(funcId+sfmt, a...)
}
