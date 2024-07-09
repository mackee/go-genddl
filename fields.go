package genddl

import (
	"fmt"
	"go/types"
	"strings"
)

type columnInfo struct {
	name       string
	typeName   string
	isNullable bool
	tagMap     map[string]string
}

var supportedTypes = map[string]struct{}{
	"time.Time": {},
}

var nullTypes = map[string]string{
	"sql.NullBool":    "bool",
	"sql.NullInt16":   "int16",
	"sql.NullInt32":   "int32",
	"sql.NullInt64":   "int64",
	"sql.NullFloat64": "float64",
	"sql.NullString":  "string",
	"sql.NullByte":    "[]byte",
	"sql.NullTime":    "time.Time",
	"sql.Null":        "",
	"mysql.NullTime":  "time.Time",
}

type columnTypeInfo struct {
	name       string
	isNullable bool
}

func toColumnTypeInfo(t types.Type) (*columnTypeInfo, error) {
	tname := namedTypeName(t)
	nt, ok := t.(*types.Named)
	if !ok {
		return &columnTypeInfo{
			name: tname,
		}, nil
	}

	if _, ok := supportedTypes[tname]; ok {
		return &columnTypeInfo{
			name:       tname,
			isNullable: false,
		}, nil
	}
	if nt, ok := nullTypes[tname]; ok {
		tname = nt
		// empty is generics types
		if nt == "" {
			nt, ok := t.(*types.Named)
			if !ok {
				return nil, fmt.Errorf("generics types must be named: %s", t.String())
			}
			tas := nt.TypeArgs()
			if tas == nil {
				return nil, fmt.Errorf("generics types must have type args: %s", t.String())
			}
			if tas.Len() != 1 {
				return nil, fmt.Errorf("generics types must have one type arg: %s", t.String())
			}
			tp := tas.At(0)
			utp, err := toColumnTypeInfo(tp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert to column type info by %s: %w", tp.String(), err)
			}
			tname = utp.name
		}
		return &columnTypeInfo{
			name:       tname,
			isNullable: true,
		}, nil
	}
	ut := nt.Underlying()
	cti, err := toColumnTypeInfo(ut)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to column type info by %s: %w", ut.String(), err)
	}
	return cti, nil
}

func namedTypeName(t types.Type) string {
	nt, ok := t.(*types.Named)
	if !ok {
		return t.String()
	}
	return strings.Join([]string{nt.Obj().Pkg().Name(), nt.Obj().Name()}, ".")
}

func columnsByFields(t types.Type, ti *types.Info, tagText string, prefix string) ([]columnInfo, error) {
	tag := parseTag(tagText)
	if tag == nil || tag["db"] == "" || tag["db"] == "-" {
		return nil, nil
	}
	if _, ok := tag["nested"]; ok {
		prefix := prefix + tag["db"]
		ret, err := columnsByStruct(t, ti, prefix)
		if err != nil {
			return nil, fmt.Errorf("failed to get columns by struct: %w", err)
		}
		return ret, nil
	}
	cti, err := toColumnTypeInfo(t)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to column type info: %w", err)
	}
	ci := columnInfo{
		name:       prefix + tag["db"],
		tagMap:     tag,
		typeName:   cti.name,
		isNullable: cti.isNullable,
	}
	return []columnInfo{ci}, nil
}

func columnsByStruct(t types.Type, ti *types.Info, prefix string) ([]columnInfo, error) {
	tn, ok := t.(*types.Named)
	if !ok {
		return nil, fmt.Errorf("type must be named: %s", t.String())
	}
	ut := tn.Underlying()
	st, ok := ut.(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("underlying type must be struct: %s", ut.String())
	}
	columns := make([]columnInfo, 0, st.NumFields())
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		tagText := st.Tag(i)
		cs, err := columnsByFields(field.Type(), ti, tagText, prefix)
		if err != nil {
			return nil, fmt.Errorf("failed to get columns by fields: %w", err)
		}
		columns = append(columns, cs...)
	}

	return columns, nil
}
