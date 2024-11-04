package parsers

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
)

var AllowedMethods = [...]string{
	http.MethodGet,
	http.MethodPut,
	http.MethodDelete,
	http.MethodPost,
	http.MethodPatch,
	http.MethodPut,
}

type HandlerParser struct {
	debugMode bool
}

type Function struct {
	FuncName   string
	FuncParams []string
	FuncPos    token.Position
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
			var isValidFunction = false

			for i := 0; i < len(AllowedMethods); i++ {
				if fn.Name.Name == AllowedMethods[i] {
					isValidFunction = true
					break
				}
			}

			if !isValidFunction {
				continue
			}

			fnParams := p.getFunctionParams(fn)

			signature.Functions = append(signature.Functions, Function{
				FuncName:   fn.Name.Name,
				FuncPos:    fs.Position(fn.Pos()),
				FuncParams: fnParams,
			})
		}
	}

	return &signature, nil
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

func (p *HandlerParser) getFunctionParams(decl *ast.FuncDecl) []string {
	var params []string = []string{}

	for _, param := range decl.Type.Params.List {
		param := p.exprToString(param.Type)

		params = append(params, param)
	}

	return params
}
