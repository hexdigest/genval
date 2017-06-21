package types

import (
	"fmt"
	"io"
)

func NewPointer(inner TypeDef) *typePointer {
	return &typePointer{
		nullable:  true,
		innerType: inner,
	}
}

type typePointer struct {
	nullable  bool
	innerType TypeDef
}

func (t typePointer) Type() string {
	return t.innerType.Type()
}

func (t *typePointer) SetValidateTag(tag Tag) error {
	switch tag.Key() {
	case PointerNullableKey:
		t.nullable = true
	case PointerNotNullKey:
		t.nullable = false
	default:
		return t.innerType.SetValidateTag(tag)
	}
	return nil
}

func (t typePointer) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.nullable {
		fmt.Fprintf(w, "if %s != nil {\n", name.Full())
		t.innerType.Generate(w, cfg, name.WithPointer())
		fmt.Fprintf(w, "}\n")
	} else {
		fmt.Fprintf(w, "if %s == nil {\n", name.Full())
		fmt.Fprintf(w, "    errs.AddFieldf(%s, \"cannot be nil\")\n", name.LabelName())
		fmt.Fprintf(w, "} else {\n")
		t.innerType.Generate(w, cfg, name.WithPointer())
		fmt.Fprintf(w, "}\n")
	}
}

func (t typePointer) Validate() error {
	return t.innerType.Validate()
}
