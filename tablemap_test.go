package genddl

import (
	"bytes"
	"go/ast"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	exampleSchemaDir       = "./_example"
	exampleExpectedSQLFile = "./_example/mysql_product.sql"
	hasIndexTableName      = "product"
)

func TestTableMap__WriteDDL(t *testing.T) {
	tr, err := retrieveTables(exampleSchemaDir)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	product := tr.tables[hasIndexTableName]
	funcs := tr.funcs[product]
	tablesRev := map[*ast.StructType]string{}
	for tableName, st := range tr.tables {
		tablesRev[st] = tableName
	}
	tm := NewTableMap(hasIndexTableName, product, funcs, tablesRev, tr.ti, false, false, false)
	bs := &bytes.Buffer{}
	err = tm.WriteDDL(bs, MysqlDialect{})
	if err != nil {
		t.Errorf("unexpected error on WriteDDL: %s", err)
	}

	ef, err := os.Open(exampleExpectedSQLFile)
	if err != nil {
		t.Fatalf("unexpected error on read expected SQL: %s", err)
	}
	_expected, err := ioutil.ReadAll(ef)
	if err != nil {
		t.Fatalf("unexpected error on read expected SQL: %s", err)
	}
	if diff := cmp.Diff(string(_expected), bs.String()); diff != "" {
		t.Errorf("result is mismatch (-want +got):\n%s", diff)
	}

}
