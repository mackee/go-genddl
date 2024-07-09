package genddl

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"log"
	"strconv"
	"strings"
)

const (
	viewSelectStatementFuncName = "_selectStatement"
)

type ViewMap struct {
	Name            string
	Columns         []string
	SelectStatement string
}

type newViewMapInput struct {
	name  string
	st    *ast.StructType
	funcs []*ast.FuncDecl
	ti    *types.Info
}

func NewViewMap(input newViewMapInput) (*ViewMap, error) {
	st := input.st
	funcs := input.funcs
	name := input.name
	columns := make([]string, 0, len(st.Fields.List))
	for _, field := range st.Fields.List {
		t := input.ti.TypeOf(field.Type)
		tagText := field.Tag.Value
		cs, err := columnsByFields(t, input.ti, tagText, "")
		if err != nil {
			return nil, fmt.Errorf("failed to get columns by fields: %w", err)
		}
		for _, c := range cs {
			columns = append(columns, c.name)
		}
	}
	selectStatement := retrieveSelectStatementByFuncs(funcs)
	if selectStatement == "" {
		return nil, nil
	}

	return &ViewMap{
		Name:            name,
		Columns:         columns,
		SelectStatement: selectStatement,
	}, nil
}

func retrieveSelectStatementByFuncs(funcs []*ast.FuncDecl) string {
	for _, funcDecl := range funcs {
		if funcDecl.Name.Name != viewSelectStatementFuncName {
			continue
		}
		var rt *ast.ReturnStmt
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			switch t := n.(type) {
			case *ast.ReturnStmt:
				rt = t
				return false
			default:
				return true
			}
		})
		if rt == nil {
			continue
		}
		for _, lit := range rt.Results {
			bl, ok := lit.(*ast.BasicLit)
			if !ok {
				continue
			}
			uq, err := strconv.Unquote(bl.Value)
			if err != nil {
				log.Printf("error by unquote: %s", err)
				return ""
			}
			uq = strings.Trim(uq, ";\n\t ")
			scanner := bufio.NewScanner(strings.NewReader(uq))
			b := &bytes.Buffer{}
			for scanner.Scan() {
				b.WriteString("  ")
				b.WriteString(scanner.Text())
				b.WriteString("\n")
			}
			return strings.TrimSuffix(b.String(), "\n")
		}
	}
	return ""
}

func (vm *ViewMap) fields(dialect Dialect) string {
	sb := &strings.Builder{}
	for i, c := range vm.Columns {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(dialect.QuoteField(c))
	}
	return sb.String()
}

func (vm *ViewMap) WriteDDL(w io.Writer, dialect Dialect) error {
	io.WriteString(w, "CREATE VIEW "+dialect.QuoteField(vm.Name)+"\n")
	io.WriteString(w, "  ("+vm.fields(dialect)+") AS\n")
	io.WriteString(w, vm.SelectStatement+";\n\n")
	return nil
}
