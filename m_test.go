package m

import (
	"testing"
)

func Test_map(t *testing.T) {
	people := []struct {
		Name  string
		State string
	}{
		{"Alice", "AZ"},
		{"Bob", "TX"},
		{"Eve", "TX"},
	}

	tests := []struct {
		Element  Element
		Expected string
	}{
		{
			M(""),
			`<div></div>`,
		},
		{
			Document(M("html[lang=en]")),
			"<!DOCTYPE html>\n" + `<html lang="en"></html>`,
		},
		{
			T(`Me & You`),
			`Me &amp; You`,
		},
		{
			S(T("A"), T("B"), T("C")),
			`ABC`,
		},
		{
			M("div", Attrf("data-x", "%d", 123)),
			`<div data-x="123"></div>`,
		},
		{
			F("Hello, %s", `Tim<`),
			`Hello, Tim&lt;`,
		},
		{
			If(true, T("true")),
			`true`,
		},
		{
			If(false, T("true")),
			``,
		},
		{
			Range(0, func(i int) Element {
				return F("%d ", i)
			}),
			``,
		},
		{
			Range(3, func(i int) Element {
				return F("%d ", i)
			}),
			`0 1 2 `,
		},
		{
			For(2, 0, -1, func(i int) Element {
				return F("%d ", i)
			}),
			`2 1 0 `,
		},
		{
			Group(len(people), func(i, j int) bool {
				return people[i].State == people[j].State
			}, func(i, j int) Element {
				return S(
					M("h2", T(people[i].State)),
					M("ul",
						For(i, j, 1, func(i int) Element {
							return M("li", T(people[i].Name))
						}),
					),
				)
			}),
			`<h2>AZ</h2><ul><li>Alice</li></ul><h2>TX</h2><ul><li>Bob</li><li>Eve</li></ul>`,
		},
		{
			Range(3, func(i int) Element {
				return M("option", If(i == 2, Attr("selected", "")), Attrf("value", "%d", i),
					F("Element %d", i),
				)
			}),
			`<option value="0">Element 0</option><option value="1">Element 1</option><option selected="" value="2">Element 2</option>`,
		},
	}

	for _, tt := range tests {
		if output := RenderString(tt.Element); output != tt.Expected {
			t.Errorf("RenderString(%#v)\ngot:\n%#v\nexpected:\n%#v", tt.Element, output, tt.Expected)
		}
	}
}
