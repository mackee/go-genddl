package genddl

import (
	"go/types"
	"strings"
)

type columnInfo struct {
	name     string
	typeName string
	tagMap   map[string]string
}

func columnsByFields(t types.Type, ti *types.Info, tagText string, prefix string) []columnInfo {
	tag := parseTag(tagText)
	if tag == nil || tag["db"] == "" || tag["db"] == "-" {
		return nil
	}
	if _, ok := tag["nested"]; ok {
		prefix := prefix + tag["db"]
		return columnsByStruct(t, ti, prefix)
	}
	name := tag["db"]
	var typeName string
	for {
		if _, ok := supportedTypes[t.String()]; ok {
			nt := t.(*types.Named)
			typeName = strings.Join([]string{nt.Obj().Pkg().Name(), nt.Obj().Name()}, ".")
			break
		}
		if ta, ok := t.(*types.Named); ok {
			t = ta.Underlying()
			continue
		} else {
			typeName = t.String()
			break
		}
	}
	return []columnInfo{
		{name: prefix + name, typeName: typeName, tagMap: tag},
	}
}

func columnsByStruct(t types.Type, ti *types.Info, prefix string) []columnInfo {
	tn, ok := t.(*types.Named)
	if !ok {
		return nil
	}
	ut := tn.Underlying()
	st, ok := ut.(*types.Struct)
	if !ok {
		return nil
	}
	columns := make([]columnInfo, 0, st.NumFields())
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		tagText := st.Tag(i)
		cs := columnsByFields(field.Type(), ti, tagText, prefix)
		columns = append(columns, cs...)
	}

	return columns
}
