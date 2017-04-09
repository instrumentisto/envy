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
				v bool `env:"BOOL"`
			}{true}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldBeFalse)
		})

		Convey("string", func() {
			setEnv("STRING", "foo")
			obj := &struct {
				v string `env:"STRING"`
			}{"bar"}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, "foo")
		})

		Convey("int", func() {
			setEnv("INT", "0")
			obj := &struct {
				v int `env:"INT"`
			}{-1}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, int(0))
		})

		Convey("int8", func() {
			setEnv("INT8", "-123")
			obj := &struct {
				v int8 `env:"INT8"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, int8(-127))
		})

		Convey("int16", func() {
			setEnv("INT16", "-32760")
			obj := &struct {
				v int16 `env:"INT16"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, int16(-32760))
		})

		Convey("int32", func() {
			setEnv("INT32", "-8388600")
			obj := &struct {
				v int32 `env:"INT32"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, int32(-8388600))
		})

		Convey("int64", func() {
			setEnv("INT64", "-2147483640")
			obj := &struct {
				v int64 `env:"INT64"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, int64(-2147483640))
		})

		Convey("uint", func() {
			setEnv("UINT", "0")
			obj := &struct {
				v uint `env:"UINT"`
			}{1}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, uint(0))
		})

		Convey("uint8", func() {
			setEnv("UINT8", "250")
			obj := &struct {
				v uint8 `env:"UINT8"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, uint8(250))
		})

		Convey("uint16", func() {
			setEnv("UINT16", "65530")
			obj := &struct {
				v uint16 `env:"UINT16"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, uint16(65530))
		})

		Convey("uint32", func() {
			setEnv("UINT32", "16777210")
			obj := &struct {
				v uint32 `env:"UINT32"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, uint32(16777210))
		})

		Convey("uint64", func() {
			setEnv("UINT64", "4294967290")
			obj := &struct {
				v uint64 `env:"UINT64"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, uint64(4294967290))
		})

		Convey("BYTE", func() {
			setEnv("BYTE", "255")
			obj := &struct {
				v byte `env:"BYTE"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, byte(255))
		})

		Convey("rune", func() {
			setEnv("RUNE", "8388600")
			obj := &struct {
				v rune `env:"RUNE"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, rune(8388600))
		})

		Convey("float32", func() {
			setEnv("FLOAT32", "3.40282346638528859811704183484516925440e+38")
			obj := &struct {
				v float32 `env:"FLOAT32"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual,
				float32(3.40282346638528859811704183484516925440e+38))
		})

		Convey("float64", func() {
			setEnv("FLOAT64", "1.797693134862315708145274237317043567981e+308")
			obj := &struct {
				v float64 `env:"FLOAT64"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual,
				float64(1.797693134862315708145274237317043567981e+308))
		})

		Convey("complex64", func() {
			setEnv("COMPLEX64", "2+3i")
			obj := &struct {
				v complex64 `env:"COMPLEX64"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, complex64(2+3i))
		})

		Convey("complex128", func() {
			setEnv("COMPLEX128", "-4352342-777774i")
			obj := &struct {
				v complex128 `env:"COMPLEX128"`
			}{}
			err := p.Parse(obj)
			So(err, ShouldBeNil)
			So(obj.v, ShouldEqual, complex128(-4352342-777774i))
		})

	})

	Convey("On unsupported type", t, func() {
		setEnv("UNSUPPORTED_TYPE", "2")
		obj := &struct {
			v uintptr `env:"UNSUPPORTED_TYPE"`
		}{5}
		err := Parser{}.Parse(obj)

		Convey("Returns error", func() {
			So(err, ShouldNotBeNil)
			// TODO: check for concrete type of error?
		})

		Convey("Does not mutate value", func() {
			So(obj.v, ShouldEqual, uintptr(5))
		})

	})

	Convey("Parses nested structs", t, func() {
		setEnv("NESTED_BOOL", "true")
		setEnv("NESTED_INT", "-10")
		setEnv("NESTED_UINT", "15")
		obj := &struct {
			v bool `env:"NESTED_BOOL"`
			n struct {
				v int `env:"NESTED_INT"`
				n struct {
					n struct {
						v uint `env:"NESTED_UINT"`
					}
				}
			}
		}{}
		err := Parser{}.Parse(obj)
		So(err, ShouldBeNil)
		So(obj.v, ShouldEqual, true)
		So(obj.n.v, ShouldEqual, int(-10))
		So(obj.n.n.n.v, ShouldEqual, uint(15))
	})

	Convey("Parses values for types behind pointers", t, func() {
		setEnv("DEREF_BOOL", "true")
		setEnv("DEREF_INT", "-10")
		v := 5; ptr1 := &v; ptr2 := &ptr1
		obj := &struct {
			v *bool `env:"DEREF_BOOL"`
			n *struct {
				v ***int `env:"DEREF_INT"`
			}
		}{n: &struct {
			v ***int `env:"DEREF_INT"`
		}{&ptr2}}
		err := Parser{}.Parse(obj)
		So(err, ShouldBeNil)
		So(*obj.v, ShouldEqual, true)
		So(***(*obj.n).v, ShouldEqual, int(-10))
	})
}

func setEnv(name, val string) {
	if err := os.Setenv(name, val); err != nil {
		panic(err.Error())
	}
}

func unsetEnv(name string) {
	if err := os.Unsetenv("CONF_DEBUG_MODE"); err != nil {
		panic(err.Error())
	}
}
