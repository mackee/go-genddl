package genddl

import (
	"bufio"
	"bytes"
	"go/ast"
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

func NewViewMap(name string, st *ast.StructType, funcs []*ast.FuncDecl) *ViewMap {
	columns := make([]string, 0, len(st.Fields.List))
	for _, field := range st.Fields.List {
		tgm := parseTag(field.Tag.Value)
		dbtag, ok := tgm["db"]
		if !ok {
			continue
		}
		columns = append(columns, dbtag)
	}
	selectStatement := retrieveSelectStatementByFuncs(funcs)
	if selectStatement == "" {
		return nil
	}

	return &ViewMap{
		Name:            name,
		Columns:         columns,
		SelectStatement: selectStatement,
	}
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

func (vm *ViewMap) WriteDDL(w io.Writer) error {
	io.WriteString(w, "CREATE VIEW "+vm.Name+" AS\n")
	io.WriteString(w, vm.SelectStatement+";\n\n")
	return nil
}
