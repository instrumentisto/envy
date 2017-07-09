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
