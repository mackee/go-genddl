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
	vm := NewViewMap("user_product", userProduct, funcs)
	if err := vm.WriteDDL(bs); err != nil {
		t.Errorf("unexpected error on WriteDDL: %s", err)
	}

	expect := `CREATE VIEW user_product AS
  SELECT p.id, u.user_id, ru.received_user_name, p.id, p.type FROM product AS p
    INNER JOIN user AS u ON p.user_id = u.id
    LEFT JOIN user AS ru ON p.received_user_id = ru.id;

`
	if diff := cmp.Diff(expect, bs.String()); diff != "" {
		t.Errorf("result is mismatch (-want +got):\n%s", diff)
	}

}
