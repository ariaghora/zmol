package val

import "github.com/ariaghora/zmol/pkg/ast"

// Function type
type ZFunction struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Env
}

func (z *ZFunction) Type() ZValueType { return "function" }
func (z *ZFunction) Str() string      { return "function" }
func (z *ZFunction) Equals(other ZValue) bool {
	if other.Type() != ZFUNCTION {
		return false
	}
	return false
}
func (z *ZFunction) LessThanEquals(other ZValue) bool {
	// For now we don't support function comparison
	return false
}
