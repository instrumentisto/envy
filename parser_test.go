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

package envigo

import (
	"errors"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser_Parse(t *testing.T) {
	Convey("If non-struct pointer is passed", t, func() {
		p := Parser{}
		obj1 := struct {
			V bool
		}{}
		some := true
		obj2 := &some

		Convey("Returns error", func() {
			So(p.Parse(obj1), ShouldEqual, ErrNotStructPtr)
			So(p.Parse(obj2), ShouldEqual, ErrNotStructPtr)
		})
	})

	Convey("Parses supported types", t, func() {
		p := Parser{}

		Convey("bool", func() {
			setEnv("BOOL", "false")
			obj := &struct {
				V bool `env:"BOOL"`
			}{true}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldBeFalse)
		})

		Convey("string", func() {
			setEnv("STRING", "foo")
			obj := &struct {
				V string `env:"STRING"`
			}{"bar"}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, "foo")
		})

		Convey("int", func() {
			setEnv("INT", "0")
			obj := &struct {
				V int `env:"INT"`
			}{-1}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, int(0))
		})

		Convey("int8", func() {
			setEnv("INT8", "-123")
			obj := &struct {
				V int8 `env:"INT8"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, int8(-123))
		})

		Convey("int16", func() {
			setEnv("INT16", "-32760")
			obj := &struct {
				V int16 `env:"INT16"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, int16(-32760))
		})

		Convey("int32", func() {
			setEnv("INT32", "-8388600")
			obj := &struct {
				V int32 `env:"INT32"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, int32(-8388600))
		})

		Convey("int64", func() {
			setEnv("INT64", "-2147483640")
			obj := &struct {
				V int64 `env:"INT64"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, int64(-2147483640))
		})

		Convey("uint", func() {
			setEnv("UINT", "0")
			obj := &struct {
				V uint `env:"UINT"`
			}{1}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, uint(0))
		})

		Convey("uint8", func() {
			setEnv("UINT8", "250")
			obj := &struct {
				V uint8 `env:"UINT8"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, uint8(250))
		})

		Convey("uint16", func() {
			setEnv("UINT16", "65530")
			obj := &struct {
				V uint16 `env:"UINT16"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, uint16(65530))
		})

		Convey("uint32", func() {
			setEnv("UINT32", "16777210")
			obj := &struct {
				V uint32 `env:"UINT32"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, uint32(16777210))
		})

		Convey("uint64", func() {
			setEnv("UINT64", "4294967290")
			obj := &struct {
				V uint64 `env:"UINT64"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, uint64(4294967290))
		})

		Convey("byte", func() {
			setEnv("BYTE", "255")
			obj := &struct {
				V byte `env:"BYTE"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, byte(255))
		})

		Convey("rune", func() {
			setEnv("RUNE", "8388600")
			obj := &struct {
				V rune `env:"RUNE"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, rune(8388600))
		})

		Convey("float32", func() {
			setEnv("FLOAT32", "3.40282346638528859811704183484516925440e+38")
			obj := &struct {
				V float32 `env:"FLOAT32"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual,
				float32(3.40282346638528859811704183484516925440e+38))
		})

		Convey("float64", func() {
			setEnv("FLOAT64", "1.797693134862315708145274237317043567981e+308")
			obj := &struct {
				V float64 `env:"FLOAT64"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual,
				float64(1.797693134862315708145274237317043567981e+308))
		})

		Convey("time.Duration", func() {
			setEnv("DURATION", "-1h2m3s4ms5us6ns")
			obj := &struct {
				V time.Duration `env:"DURATION"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual,
				-(time.Hour + 2*time.Minute + 3*time.Second +
					4*time.Millisecond + 5*time.Microsecond +
					6*time.Nanosecond))
		})

	})

	Convey("On unsupported type", t, func() {
		setEnv("UNSUPPORTED_TYPE", "2")
		obj := &struct {
			V uintptr `env:"UNSUPPORTED_TYPE"`
		}{5}
		err := Parser{}.Parse(obj)

		Convey("Returns error", func() {
			So(err, ShouldNotBeNil)
			So(err, ShouldHaveSameTypeAs, UnparsableTypeError{})
		})

		Convey("Does not mutate value", func() {
			So(obj.V, ShouldEqual, uintptr(5))
		})
	})

	Convey("On incorrectly declared tag", t, func() {
		setEnv("UINT8", "3")
		obj := &struct {
			V uint8 `env:""`
		}{5}
		err := Parser{}.Parse(obj)

		Convey("Returns error", func() {
			So(err, ShouldNotBeNil)
			So(err, ShouldHaveSameTypeAs, EmptyVarNameError{})
		})

		Convey("Does not mutate value", func() {
			So(obj.V, ShouldEqual, 5)
		})
	})

	Convey("Parses nested structs", t, func() {
		setEnv("NESTED_BOOL", "true")
		setEnv("NESTED_INT", "-10")
		setEnv("NESTED_UINT", "15")
		obj := &struct {
			V bool `env:"NESTED_BOOL"`
			N struct {
				V int `env:"NESTED_INT"`
				N struct {
					N struct {
						V uint `env:"NESTED_UINT"`
					}
				}
			}
			h bool // nolint: unused, megacheck
		}{}
		err := Parser{}.Parse(obj)

		So(err, ShouldBeNil)
		So(obj.V, ShouldEqual, true)
		So(obj.N.V, ShouldEqual, int(-10))
		So(obj.N.N.N.V, ShouldEqual, uint(15))
	})

	Convey("Parses embedded structs", t, func() {
		setEnv("EMBEDDED_BOOL", "true")
		setEnv("EMBEDDED_INT", "-2")

		Convey("Raw struct", func() {
			obj := &struct {
				EmbeddedStruct
				V bool `env:"EMBEDDED_BOOL"`
			}{}
			err := Parser{}.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, true)
			So(obj.EmbeddedStruct.V, ShouldEqual, true)
			So(obj.EmbeddedStruct.V2, ShouldEqual, -2)
		})

		Convey("Struct behind pointer", func() {
			obj := &struct {
				*EmbeddedStruct
				V bool `env:"EMBEDDED_BOOL"`
			}{EmbeddedStruct: &EmbeddedStruct{}}
			err := Parser{}.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, true)
			So(obj.EmbeddedStruct.V, ShouldEqual, true)
			So(obj.EmbeddedStruct.V2, ShouldEqual, -2)
		})
	})

	Convey("Parses values for types behind pointers", t, func() {
		setEnv("DEREF_BOOL", "true")
		setEnv("DEREF_INT", "-10")
		i := 5
		ptr1 := &i
		ptr2 := &ptr1
		b := true
		obj := &struct {
			V *bool `env:"DEREF_BOOL"`
			N *struct {
				V ***int `env:"DEREF_INT"`
			}
		}{V: &b, N: &struct {
			V ***int `env:"DEREF_INT"`
		}{&ptr2}}
		err := Parser{}.Parse(obj)

		So(err, ShouldBeNil)
		So(*(obj.V), ShouldEqual, true)
		So(***(obj.N).V, ShouldEqual, int(-10))
	})

	Convey("Omitts nil pointers", t, func() { // TODO: do not omit nil pointers!
		setEnv("PTR_BOOL", "true")
		obj := &struct {
			V *bool `env:"PTR_BOOL"`
		}{}
		err := Parser{}.Parse(obj)

		So(err, ShouldBeNil)
		So(obj.V, ShouldBeNil)
	})

	Convey("Uses custom parser if type has one", t, func() {

		Convey("Performs custome parse correctly", func() {
			setEnv("CUSTOM_UINT8", "10")
			v1 := customUint8(1)
			v2 := customUint8(2)
			pv2 := &v2
			obj1 := &struct {
				V *customUint8 `env:"CUSTOM_UINT8"`
			}{&v1}
			obj2 := &struct {
				V **customUint8 `env:"CUSTOM_UINT8"`
			}{&pv2}
			err1 := Parser{}.Parse(obj1)
			err2 := Parser{}.Parse(obj2)

			So(err1, ShouldBeNil)
			So(*(obj1.V), ShouldEqual, 7)
			So(err2, ShouldBeNil)
			So(**(obj2.V), ShouldEqual, 7)
		})

	})

	Convey("On tagged struct without custom parser", t, func() {

		Convey("Returns error of unsupported type", func() {
			setEnv("EMBEDDED_STRUCT", "")
			setEnv("EMBEDDED_BOOL", "true")
			setEnv("EMBEDDED_INT", "-2")
			obj := &struct {
				V EmbeddedStruct `env:"EMBEDDED_STRUCT"`
			}{}
			err := Parser{}.Parse(obj)

			So(err, ShouldNotBeNil)
			So(err, ShouldHaveSameTypeAs, UnparsableTypeError{})
		})

	})

	Convey("If value cannot be parsed", t, func() {

		Convey("Returns parsing error", func() {

			Convey("supported types", func() {
				setEnv("FAIL_BOOL", "hi")
				setEnv("FAIL_INT", "true")
				setEnv("FAIL_UINT", "false")
				setEnv("FAIL_FLOAT", "-----")
				setEnv("FAIL_DURATION", "???")
				obj1 := &struct {
					V bool `env:"FAIL_BOOL"`
				}{}
				obj2 := &struct {
					V int `env:"FAIL_INT"`
				}{}
				obj3 := &struct {
					V uint16 `env:"FAIL_UINT"`
				}{}
				obj4 := &struct {
					V float64 `env:"FAIL_FLOAT"`
				}{}
				obj5 := &struct {
					V time.Duration `env:"FAIL_DURATION"`
				}{}
				err1 := Parser{}.Parse(obj1)
				err2 := Parser{}.Parse(obj2)
				err3 := Parser{}.Parse(obj3)
				err4 := Parser{}.Parse(obj4)
				err5 := Parser{}.Parse(obj5)

				So(err1, ShouldNotBeNil)
				So(err1, ShouldHaveSameTypeAs, ParseError{})
				So(err1.Error(), ShouldContainSubstring, "'FAIL_BOOL'")
				So(err2, ShouldNotBeNil)
				So(err2, ShouldHaveSameTypeAs, ParseError{})
				So(err2.Error(), ShouldContainSubstring, "'FAIL_INT'")
				So(err3, ShouldNotBeNil)
				So(err3, ShouldHaveSameTypeAs, ParseError{})
				So(err3.Error(), ShouldContainSubstring, "'FAIL_UINT'")
				So(err4, ShouldNotBeNil)
				So(err4, ShouldHaveSameTypeAs, ParseError{})
				So(err4.Error(), ShouldContainSubstring, "'FAIL_FLOAT'")
				So(err5, ShouldNotBeNil)
				So(err5, ShouldHaveSameTypeAs, ParseError{})
				So(err5.Error(), ShouldContainSubstring, "'FAIL_DURATION'")
			})

			Convey("custom parser type", func() {
				setEnv("FAIL_CUSTOM", "10")
				v := customFailure(4)
				pV := &v
				obj1 := &struct {
					V *customFailure `env:"FAIL_CUSTOM"`
				}{pV}
				obj2 := &struct {
					V **customFailure `env:"FAIL_CUSTOM"`
				}{&pV}
				err1 := Parser{}.Parse(obj1)
				err2 := Parser{}.Parse(obj2)

				So(err1, ShouldNotBeNil)
				So(err1, ShouldHaveSameTypeAs, ParseError{})
				So(err2, ShouldNotBeNil)
				So(err2, ShouldHaveSameTypeAs, ParseError{})
			})

			Convey("nested structs", func() {
				setEnv("FAIL_NESTED", "?-?")
				obj1 := &struct {
					N struct {
						V int `env:"FAIL_NESTED"`
					}
				}{}
				obj2 := &struct {
					N *struct {
						V int `env:"FAIL_NESTED"`
					}
				}{&struct {
					V int `env:"FAIL_NESTED"`
				}{}}
				err1 := Parser{}.Parse(obj1)
				err2 := Parser{}.Parse(obj2)

				So(err1, ShouldNotBeNil)
				So(err1, ShouldHaveSameTypeAs, ParseError{})
				So(err2, ShouldNotBeNil)
				So(err2, ShouldHaveSameTypeAs, ParseError{})
			})

		})

	})
}

type customUint8 uint8

func (v *customUint8) UnmarshalText(_ []byte) error {
	*v = 7
	return nil
}

type customFailure uint8

func (_ *customFailure) UnmarshalText(_ []byte) error {
	return errors.New("some error")
}

type EmbeddedStruct struct {
	V  bool `env:"EMBEDDED_BOOL"`
	V2 int  `env:"EMBEDDED_INT"`
}

// setEnv is a simple helper function for setting env vars in one line.
func setEnv(name, val string) {
	if err := os.Setenv(name, val); err != nil {
		panic(err)
	}
}
