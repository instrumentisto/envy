package envigo

import (
	"encoding"
	"errors"
	"os"
	refl "reflect"
	"strconv"
	"time"
)

var ErrUnsupportedType = errors.New("envifo: unsupported type")

type Parser struct{}

func (p Parser) Parse(obj interface{}) error {
	// todo: not a struct pointer-check
	return p.parseStruct(refl.ValueOf(obj).Elem())
}

func (p Parser) parseStruct(structVal refl.Value) error {
	structType := structVal.Type()
L:
	for i := 0; i < structVal.NumField(); i++ {
		fieldVal := structVal.Field(i)

		// Omit private field
		if !fieldVal.CanSet() {
			continue
		}

		// Dereference pointer
		for fieldVal.Kind() == refl.Ptr {
			if fieldVal.IsNil() {
				continue L
			}
			fieldVal = fieldVal.Elem()
		}
		fieldKind := fieldVal.Kind()

		// Omit non-struct or parse recursive struct if no `env` tag
		envVarName, hasTag := structType.Field(i).Tag.Lookup("env")
		if !hasTag {
			if fieldKind != refl.Struct {
				continue
			} else {
				if err := p.parseStruct(fieldVal); err != nil {
					return err
				}
			}
		}
		envValue := os.Getenv(envVarName) // TODO: do not parse always

		// Unmarshal with custom unmarshaller
		if field, ok := fieldVal.Interface().(encoding.TextUnmarshaler); ok {
			if err := field.UnmarshalText([]byte(envValue)); err != nil {
				return err
			}
			continue
		}
		// Unmarshal as time.Duration
		fieldType := fieldVal.Type()
		if fieldType.PkgPath() == "time" && fieldType.Name() == "Duration" {
			val, err := time.ParseDuration(envValue)
			if err != nil {
				return err
			}
			fieldVal.SetInt(int64(val))
			continue
		}
		// Unmarshal as primitive type
		switch fieldKind {
		case refl.Bool:
			val, err := strconv.ParseBool(envValue)
			if err != nil {
				return err
			}
			fieldVal.SetBool(val)
		case refl.String:
			fieldVal.SetString(envValue)
		case refl.Int, refl.Int8, refl.Int16, refl.Int32, refl.Int64:
			val, err := strconv.ParseInt(envValue, 0, fieldVal.Type().Bits())
			if err != nil {
				return err
			}
			fieldVal.SetInt(val)
		case refl.Uint, refl.Uint8, refl.Uint16, refl.Uint32, refl.Uint64:
			val, err := strconv.ParseUint(envValue, 0, fieldVal.Type().Bits())
			if err != nil {
				return err
			}
			fieldVal.SetUint(val)
		case refl.Float32, refl.Float64:
			val, err := strconv.ParseFloat(envValue, fieldVal.Type().Bits())
			if err != nil {
				return err
			}
			fieldVal.SetFloat(val)
		// Parse recursive struct
		case refl.Struct:
			if err := p.parseStruct(fieldVal); err != nil {
				return err
			}
		default:
			return ErrUnsupportedType
		}
	}
	return nil
}
