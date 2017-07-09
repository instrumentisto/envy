package envigo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEmptyVarNameError_Error(t *testing.T) {
	Convey("Contains struct field name", t, func() {
		err := EmptyVarNameError{"field"}

		So(err.Error(), ShouldContainSubstring, "'field'")
	})
}

func TestUnparsableTypeError_Error(t *testing.T) {
	Convey("Contains struct field name", t, func() {
		err := UnparsableTypeError{"fld"}

		So(err.Error(), ShouldContainSubstring, "'fld'")
	})
}

func TestParseError_Error(t *testing.T) {
	Convey("Contains struct field name", t, func() {
		err := ParseError{"f1eld", "", ""}

		So(err.Error(), ShouldContainSubstring, "'f1eld'")
	})

	Convey("Contains env var name", t, func() {
		err := ParseError{"", "ENV_VAR", ""}

		So(err.Error(), ShouldContainSubstring, "'ENV_VAR'")
	})

	Convey("Contains error reason", t, func() {
		err := ParseError{"", "", "some reason here"}

		So(err.Error(), ShouldContainSubstring, "some reason here")
	})
}
