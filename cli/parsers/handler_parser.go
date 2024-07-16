package parsers

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type HandlerParser struct {
	debugMode bool
}

type Function struct {
	FunctionName    string
	FunctionsParams []string
	FunctionPos     token.Position
}

type HandlerSignature struct {
	PackageName string
	Functions   []Function
}

func (p *HandlerParser) Parse(file_path string) (*HandlerSignature, error) {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, file_path, nil, parser.AllErrors)

	if err != nil {
		return nil, err
	}

	var signature HandlerSignature = HandlerSignature{
		PackageName: node.Name.Name,
	}

	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			fnParams := p.getFunctionParams(fn)

			signature.Functions = append(signature.Functions, Function{
				FunctionName:    fn.Name.Name,
				FunctionPos:     fs.Position(fn.Pos()),
				FunctionsParams: fnParams,
			})
		}
	}

	return &signature, nil
}

func (p *HandlerParser) getFunctionParams(decl *ast.FuncDecl) []string {
	var params []string = []string{}

	for _, param := range decl.Type.Params.List {
		param := p.exprToString(param.Type)

		params = append(params, param)
	}

	return params
}

func (p *HandlerParser) exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		return p.exprToString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + p.exprToString(e.X)
	case *ast.Ident:
		return e.Name
	}

	return ""
}
