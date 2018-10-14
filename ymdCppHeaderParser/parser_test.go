package ymdCppHeaderParser

import (
	"testing"
)

const content = `
// Run with: header-parser example1.h -c TCLASS -e TENUM -f TFUNC -p TPROPERTY

#include <vector>

namespace test 
{
	class Foo : public Bar 
	{
	protected:
		bool ProtectedFunction(std::vector<int> args) const;

	public:
		enum Enum
    {
      FirstValue,
      SecondValue = 3
    };

		virtual void inProgress() = 0;

	public:
		int ThisIsAProperty;
  };
}`

func TestParseFile(t *testing.T) {
	p := NewParser([]byte(content))
	p.ParseAll()
}

func TestParser_ParseTypeNode(t *testing.T) {
	p := NewParser([]byte(`void test1(int , int);`))
	node := p.ParseTypeNode()
	assert(node.NodeType == kLiteral)
	assert(node.LiteralName == `void`)
	assert(p.cursorPos == 5, p.cursorPos)
}

func TestParser_ParseAccessControl(t *testing.T) {
	p := NewParser([]byte(`protected:`))
	var token Token
	ok := p.GetToken(&token, false, false)
	assert(ok)
	var access AccessControlType
	ok = p.ParseAccessControl(&token, &access)
	assert(ok && access == kProtected)
	assert(p.MatchSymbol(`:`))
	assert(p.cursorPos == 10)
}

func TestParser_ParseFunction(t *testing.T) {
	p := NewParser([]byte(`user::bool ProtectedFunction(std::map<int, std::string> ) const ;`))
	ok := p.ParseFunction()
	assert(ok)
}

func TestParser_ParseClass(t *testing.T) {
	p := NewParser([]byte(`class A : public BaseA { protected : virtual void function1() ; };`))
	p.ParseClass()
}

func TestParser_ParseAll(t *testing.T) {
	p := NewParser([]byte(`class A : public BaseA { protected : virtual void function1() ; };`))
	p.ParseAll()
}

func TestParser_ParseEnum(t *testing.T) {
	p := NewParser([]byte(`enum class UserType { UT_Student = "value1", UT_Teacher, UT_Admin, }; `))
	p.ParseEnum()
}

func TestParser_ParseDeclaration(t *testing.T) {
	p := NewParser([]byte(`int a;`))
	var token Token
	ok1 := p.GetToken(&token, false, false)
	assert(ok1)
	ok2 := p.ParseDeclaration(&token)
	assert(ok2)
}

func TestParser_ParseDirective(t *testing.T) {
	p := NewParser([]byte(`#include <stdio.h>`))
	p.ParseDirective()
}

