package config

import (
	"fmt"
	"reflect"
	"strings"
)

func validate(cfg any) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		if optional := field.Tag.Get("optional"); optional == "true" {
			continue
		}

		value := v.Field(i)

		var isEmpty bool

		switch value.Kind() { //nolint:exhaustive
		case reflect.String:
			isEmpty = strings.TrimSpace(value.String()) == ""
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			isEmpty = value.Int() == 0
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			isEmpty = value.Uint() == 0
		case reflect.Float32, reflect.Float64:
			isEmpty = value.Float() == 0
		case reflect.Ptr, reflect.Interface:
			isEmpty = value.IsNil()
		case reflect.Slice, reflect.Map, reflect.Array:
			isEmpty = value.Len() == 0
		case reflect.Struct:
			if err := validate(value.Addr().Interface()); err != nil {
				return fmt.Errorf("%s.%s: %w", t.Name(), field.Name, err)
			}
		// reflect.Bool, reflect.Uintptr, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Invalid
		default:
		}

		if isEmpty {
			return fmt.Errorf("config field [%s] is required but empty", field.Name)
		}
	}

	return nil
}

func merge(dst, src any) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	for i := range srcVal.NumField() {
		srcField := srcVal.Field(i)
		if !srcField.IsZero() {
			dstVal.Field(i).Set(srcField)
		}
	}
}
