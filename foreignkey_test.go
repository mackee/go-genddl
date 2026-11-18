package genddl

import (
	"bytes"
	"go/ast"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// tableDDL returns the CREATE TABLE output of one table followed by the
// statements that WriteDDL defers to the end of the DDL file.
func tableDDL(t *testing.T, dir, name string, opts *tableMapOption) (string, error) {
	t.Helper()
	tr, err := retrieveTables(dir)
	if err != nil {
		t.Fatalf("unexpected error on retrieveTables: %s", err)
	}
	tables := map[*ast.StructType]string{}
	for tableName, st := range tr.tables {
		tables[st] = tableName
	}
	st := tr.tables[name]
	tm, err := NewTableMap(name, st, tr.funcs[st], tables, tr.ti, opts)
	if err != nil {
		return "", err
	}
	bs := &bytes.Buffer{}
	dialect := PostgresqlDialect{}
	if err := tm.WriteDDL(bs, dialect, opts); err != nil {
		t.Fatalf("unexpected error on WriteDDL: %s", err)
	}
	for _, idx := range tm.EndOfDDLFileIndexes {
		bs.WriteString(idx.Index(dialect, tables) + ";\n")
	}
	return bs.String(), nil
}

func TestForeignKey__Composite(t *testing.T) {
	got, err := tableDDL(t, "./testdata/compositefk", "a_item_review", &tableMapOption{withoutDropTable: true})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := `CREATE TABLE "a_item_review" (
    "id" BIGINT NOT NULL,
    "shop_id" BIGINT NOT NULL,
    "item_id" BIGINT NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("shop_id", "item_id") REFERENCES "items"("shop_id", "id") ON DELETE NO ACTION
) ;

`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("result is mismatch (-want +got):\n%s", diff)
	}
}

func TestForeignKey__CompositeOuter(t *testing.T) {
	got, err := tableDDL(t, "./testdata/compositefk", "a_item_review", &tableMapOption{withoutDropTable: true, outerForeignKey: true})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := `CREATE TABLE "a_item_review" (
    "id" BIGINT NOT NULL,
    "shop_id" BIGINT NOT NULL,
    "item_id" BIGINT NOT NULL,
    PRIMARY KEY ("id")
) ;

ALTER TABLE "a_item_review" ADD CONSTRAINT "fk_a_item_review_shop_id_item_id" FOREIGN KEY ("shop_id", "item_id") REFERENCES "items"("shop_id", "id") ON DELETE NO ACTION;
`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("result is mismatch (-want +got):\n%s", diff)
	}
}

func TestForeignKey__SelfReference(t *testing.T) {
	got, err := tableDDL(t, "./testdata/selfref", "category", &tableMapOption{withoutDropTable: true})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := `CREATE TABLE "category" (
    "id" BIGINT NOT NULL PRIMARY KEY,
    "parent_id" BIGINT NOT NULL,
    FOREIGN KEY ("parent_id") REFERENCES "category"("id") ON DELETE CASCADE
) ;

`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("result is mismatch (-want +got):\n%s", diff)
	}
}

func TestForeignKey__Invalid(t *testing.T) {
	cases := []struct {
		table string
		want  string
	}{
		{"count_mismatch", "table count_mismatch: foreign key has 2 columns but 1 reference columns"},
		{"no_columns", "table no_columns: foreign key has no columns"},
		{"unresolved", "table unresolved: foreign key has a column that does not resolve to a field of a table struct"},
		{"span_structs", "table span_structs: foreign key reference columns span more than one struct"},
	}
	for _, c := range cases {
		t.Run(c.table, func(t *testing.T) {
			_, err := tableDDL(t, "./testdata/invalidfk", c.table, &tableMapOption{})
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("error %q does not contain %q", err, c.want)
			}
		})
	}
}
