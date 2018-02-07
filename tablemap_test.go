package genddl

import (
	"bytes"
	"go/ast"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	exampleSchemaDir       = "./_example"
	exampleExpectedSQLFile = "./_example/mysql_product.sql"
	hasIndexTableName      = "product"
)

func TestTableMap__WriteDDL(t *testing.T) {
	tables, funcMap, err := retrieveTables(exampleSchemaDir)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	product := tables[hasIndexTableName]
	funcs := funcMap[product]
	tablesRev := map[*ast.StructType]string{}
	for tableName, st := range tables {
		tablesRev[st] = tableName
	}
	tm := NewTableMap(hasIndexTableName, product, funcs, tablesRev)
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
	if bs.String() != string(_expected) {
		t.Errorf("result is not match: \n%s\n\t\tvs\n%s", bs.String(), string(_expected))
	}

}
