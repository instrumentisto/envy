package envy

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
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

		Convey("BYTE", func() {
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

	})

	Convey("On unsupported type", t, func() {
		setEnv("UNSUPPORTED_TYPE", "2")
		obj := &struct {
			V uintptr `env:"UNSUPPORTED_TYPE"`
		}{5}
		err := Parser{}.Parse(obj)

		Convey("Returns error", func() {
			So(err, ShouldNotBeNil)
			// TODO: check for concrete type of error?
		})

		Convey("Does not mutate value", func() {
			So(obj.V, ShouldEqual, uintptr(5))
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
		}{}
		err := Parser{}.Parse(obj)
		So(err, ShouldBeNil)
		So(obj.V, ShouldEqual, true)
		So(obj.N.V, ShouldEqual, int(-10))
		So(obj.N.N.N.V, ShouldEqual, uint(15))
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
}

func setEnv(name, val string) {
	if err := os.Setenv(name, val); err != nil {
		panic(err.Error())
	}
}

func unsetEnv(name string) {
	if err := os.Unsetenv(name); err != nil {
		panic(err.Error())
	}
}
