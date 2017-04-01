// Copyright (c) 2017 Ernest Micklei
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package proto

import "testing"

func TestOptionCases(t *testing.T) {
	for i, each := range []struct {
		proto     string
		name      string
		strLit    string
		nonStrLit string
	}{{
		`option (full).java_package = "com.example.foo";`,
		"(full).java_package",
		"com.example.foo",
		"",
	}, {
		`option Bool = true;`,
		"Bool",
		"",
		"true",
	}, {
		`option Float = -3.14E1;`,
		"Float",
		"",
		"-3.14E1",
	}, {
		`option (foo_options) = { opt1: 123 opt2: "baz" };`,
		"(foo_options)",
		"",
		"",
	}} {
		p := newParserOn(each.proto)
		pr, err := p.Parse()
		if err != nil {
			t.Fatal("testcase failed:", i, err)
		}
		if got, want := len(pr.Elements), 1; got != want {
			t.Fatalf("[%d] got [%v] want [%v]", i, got, want)
		}
		o := pr.Elements[0].(*Option)
		if got, want := o.Name, each.name; got != want {
			t.Errorf("[%d] got [%v] want [%v]", i, got, want)
		}
		if len(each.strLit) > 0 {
			if got, want := o.Constant.Source, each.strLit; got != want {
				t.Errorf("[%d] got [%v] want [%v]", i, got, want)
			}
		}
		if len(each.nonStrLit) > 0 {
			if got, want := o.Constant.Source, each.nonStrLit; got != want {
				t.Errorf("[%d] got [%v] want [%v]", i, got, want)
			}
		}
		if got, want := o.IsEmbedded, false; got != want {
			t.Errorf("[%d] got [%v] want [%v]", i, got, want)
		}
	}
}

func TestLiteralString(t *testing.T) {
	proto := `"string"`
	p := newParserOn(proto)
	l := new(Literal)
	if err := l.parse(p); err != nil {
		t.Fatal(err)
	}
	if got, want := l.IsString, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := l.Source, "string"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
