package genddl

import (
	"bytes"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestViewMap(t *testing.T) {
	tr, err := retrieveTables(exampleSchemaDir)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	userProduct := tr.views["user_product"]
	funcs := tr.funcs[userProduct]

	bs := &bytes.Buffer{}
	vm := NewViewMap(newViewMapInput{
		name:  "user_product",
		st:    userProduct,
		funcs: funcs,
		ti:    tr.ti,
	})
	if err := vm.WriteDDL(bs, MysqlDialect{}); err != nil {
		t.Errorf("unexpected error on WriteDDL: %s", err)
	}

	expect :=
		"CREATE VIEW `user_product`\n" +
			"  (`p_id`, `u_name`, `ru_name`, `p_id`, `p_type`) AS\n" +
			"  SELECT p.id, u.name, ru.name, p.id, p.type FROM product AS p\n" +
			"    INNER JOIN user AS u ON p.user_id = u.id\n" +
			"    LEFT JOIN user AS ru ON p.received_user_id = ru.id;\n\n"

	if diff := cmp.Diff(expect, bs.String()); diff != "" {
		t.Errorf("result is mismatch (-want +got):\n%s", diff)
	}

}
