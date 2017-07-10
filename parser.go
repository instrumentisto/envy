// Copyright 2017 tyranron
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package envigo implements parser to parse values from environment variables
// into tagged struct fields.
package envigo

import (
	"encoding"
	"os"
	refl "reflect"
	"strconv"
	"time"
)

// TODO: think about different behavior/mode

// Parser is an implementation of environment variables parser.
type Parser struct{}

// Parse inspects given struct and parses environment variables that were
// mentioned in struct field tag `env`.
func (p Parser) Parse(obj interface{}) error {
	ptr := refl.ValueOf(obj)
	if ptr.Kind() != refl.Ptr {
		return ErrNotStructPtr
	}
	val := ptr.Elem()
	if val.Kind() != refl.Struct {
		return ErrNotStructPtr
	}
	return p.parseStruct(val)
}

// parseStruct performs parsing for given struct.
func (p Parser) parseStruct(structVal refl.Value) error { // nolint: gocyclo
	structType := structVal.Type()
L:
	for i := 0; i < structVal.NumField(); i++ {
		fieldVal := structVal.Field(i)

		// Omit private field
		if !fieldVal.CanSet() {
			continue
		}

		envName, hasTag := structType.Field(i).Tag.Lookup("env")
		if hasTag && envName == "" {
			return EmptyVarNameError{structType.Field(i).Name}
		}
		envValue := ""
		if hasTag {
			if envValue = os.Getenv(envName); envValue == "" {
				continue
			}
		}

		// Unmarshal with custom unmarshaller
		if hasTag {
			if ok, err := parseAsTextUnmarshaler(fieldVal, envValue); ok {
				if err != nil {
					return ParseError{
						structType.Field(i).Name, envName, err.Error(),
					}
				}
				continue
			}
		}

		// Dereference pointer
		for fieldVal.Kind() == refl.Ptr {
			if fieldVal.IsNil() {
				continue L
			}
			fieldVal = fieldVal.Elem()
			if hasTag {
				if ok, err := parseAsTextUnmarshaler(fieldVal, envValue); ok {
					if err != nil {
						return ParseError{
							structType.Field(i).Name, envName, err.Error(),
						}
					}
					continue L
				}
			}
		}
		fieldKind := fieldVal.Kind()

		// If no `env` tag: omit and parse recursively if struct
		if !hasTag {
			if fieldKind == refl.Struct {
				if err := p.parseStruct(fieldVal); err != nil {
					return ParseError{
						structType.Field(i).Name, envName, err.Error(),
					}
				}
			}
			continue
		}

		// Unmarshal as time.Duration
		fieldType := fieldVal.Type()
		if fieldType.PkgPath() == "time" && fieldType.Name() == "Duration" {
			val, err := time.ParseDuration(envValue)
			if err != nil {
				return ParseError{
					structType.Field(i).Name, envName, err.Error(),
				}
			}
			fieldVal.SetInt(int64(val))
			continue
		}
		// Unmarshal as primitive type
		switch fieldKind {
		case refl.Bool:
			val, err := strconv.ParseBool(envValue)
			if err != nil {
				return ParseError{
					structType.Field(i).Name, envName, err.Error(),
				}
			}
			fieldVal.SetBool(val)
		case refl.String:
			fieldVal.SetString(envValue)
		case refl.Int, refl.Int8, refl.Int16, refl.Int32, refl.Int64:
			val, err := strconv.ParseInt(envValue, 0, fieldVal.Type().Bits())
			if err != nil {
				return ParseError{
					structType.Field(i).Name, envName, err.Error(),
				}
			}
			fieldVal.SetInt(val)
		case refl.Uint, refl.Uint8, refl.Uint16, refl.Uint32, refl.Uint64:
			val, err := strconv.ParseUint(envValue, 0, fieldVal.Type().Bits())
			if err != nil {
				return ParseError{
					structType.Field(i).Name, envName, err.Error(),
				}
			}
			fieldVal.SetUint(val)
		case refl.Float32, refl.Float64:
			val, err := strconv.ParseFloat(envValue, fieldVal.Type().Bits())
			if err != nil {
				return ParseError{
					structType.Field(i).Name, envName, err.Error(),
				}
			}
			fieldVal.SetFloat(val)
		default:
			return UnparsableTypeError{structType.Field(i).Name}
		}
	}
	return nil
}

// parseAsTextUnmarshaler tries to parse value from environment variable
// with encoding.TextUnmarshaler implementation.
func parseAsTextUnmarshaler(
	fieldVal refl.Value, envValue string,
) (bool, error) {
	if field, ok := fieldVal.Interface().(encoding.TextUnmarshaler); ok {
		return true, field.UnmarshalText([]byte(envValue))
	}
	return false, nil
}
