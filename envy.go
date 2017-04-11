package envy

import (
	"errors"
	"os"
	refl "reflect"
	"strconv"
)

var (
	ErrUnsupportedType = errors.New("envy: unsupprted type") // TODO: more info
)

type Parser struct{}

func (p Parser) Parse(obj interface{}) error {
	return p.parseStruct(refl.ValueOf(obj).Elem())
}

func (p Parser) parseStruct(structVal refl.Value) error {
	structType := structVal.Type()
	for i := 0; i < structVal.NumField(); i++ {
		fieldVal := structVal.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		for fieldVal.Kind() == refl.Ptr {
			fieldVal = fieldVal.Elem()
		}
		fieldKind := fieldVal.Kind()
		envVarName, hasTag := structType.Field(i).Tag.Lookup("env")
		if !hasTag && (fieldKind != refl.Struct) {
			continue
		}
		envValue := os.Getenv(envVarName) // TODO: do not parse always
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
