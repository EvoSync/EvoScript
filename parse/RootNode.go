package parse

import "EvoScript/lexer"

/*
	NodeType[1]: Declaration
	NodeType[2]: Return statement
	NodeType[3]: Variable expression
	NodeType[4]: Function expression
	NodeType[5]: Unknown laid
	NodeType[6]: Text
	NodeType[7]: If statement
	NodeType[8]: Function creation
*/

type Instruction struct {
	NodeType               int
	TokensAxis             []lexer.Token
	DeclarationInstruction *DeclareStatement
	FunctionInstruction    *FunctionCall
	Text                   *lexer.Token
	IF                     *IfStatement
	FunctionCreation       *FunctionBody
}
