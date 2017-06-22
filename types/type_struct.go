package types

import (
	"fmt"
	"io"
	"strings"
)

func NewStruct(typeName string) *typeStruct {
	return &typeStruct{typeName: typeName, external: false}
}
func NewExternalStruct(typeName string) *typeStruct {
	return &typeStruct{typeName: typeName, external: true}
}

type typeStruct struct {
	typeName string
	external bool
	funcs    []string
}

func (t typeStruct) Type() string {
	return t.typeName
}

func (t *typeStruct) SetValidateTag(tag Tag) error {
	switch tag.Key() {
	case StructFuncKey:
		for _, v := range parseFuncsParam(tag.(SimpleTag).Param) {
			t.funcs = append(t.funcs, v)
		}
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeStruct) Generate(w io.Writer, cfg GenConfig, name Name) {
	registerError := `errs.AddField(%s, err)`

	if !cfg.SeveralErrors {
		cfg.AddImport("fmt")
		registerError = "	return fmt.Errorf(\"%%s %%v\", %s, err)\n"
	}

	switch {
	case len(t.funcs) != 0:
		for _, f := range t.funcs {
			if strings.HasPrefix(f, ".") {
				fmt.Fprintf(w, "if err := %s%s(); err != nil {\n", name.WithoutPointer(), f)
				fmt.Fprintf(w, registerError, name.LabelName())
				fmt.Fprintf(w, "}\n")
			} else {
				fmt.Fprintf(w, "if err := %s(%s); err != nil {\n", f, name.Full())
				fmt.Fprintf(w, registerError, name.LabelName())
				fmt.Fprintf(w, "}\n")
			}
		}
	case !cfg.NeedValidatableCheck, !t.external:
		fmt.Fprintf(w, "if err := %s.Validate%s(); err != nil {\n", name.WithAlias(), name.tagName)
		fmt.Fprintf(w, registerError, name.LabelName())
		fmt.Fprintf(w, "}\n")
	default:
		fmt.Fprintf(w, "if err := validate%s(%s); err != nil {\n", name.tagName, name.WithAlias())
		fmt.Fprintf(w, registerError, name.LabelName())
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeStruct) Validate() error {
	return nil
}
