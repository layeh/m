package selector

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		Input  string
		Error  bool
		Result Result
	}{
		{
			"",
			false,
			Result{},
		},
		{
			"div",
			false,
			Result{
				TagName: "div",
			},
		},
		{
			"#hello",
			false,
			Result{
				TagName: "",
				ID:      "hello",
			},
		},
		{
			Input: "a#main.a.b",
			Error: false,
			Result: Result{
				TagName: "a",
				ID:      "main",
				Classes: []string{"a", "b"},
			},
		},
		{
			Input: ".a.b",
			Error: false,
			Result: Result{
				Classes: []string{"a", "b"},
			},
		},
		{
			Input: ".fff",
			Error: false,
			Result: Result{
				Classes: []string{"fff"},
			},
		},
		{
			Input: "[title=Hello]",
			Error: false,
			Result: Result{
				Attributes: [][2]string{
					{"title", "Hello"},
				},
			},
		},
		{
			Input: "option[selected]",
			Error: false,
			Result: Result{
				TagName: "option",
				Attributes: [][2]string{
					{"selected", ""},
				},
			},
		},
		{
			Input: "div[title=Hello][id=main]",
			Error: false,
			Result: Result{
				TagName: "div",
				Attributes: [][2]string{
					{"title", "Hello"},
					{"id", "main"},
				},
			},
		},
		{
			Input: "P.active[title=Open]",
			Error: false,
			Result: Result{
				TagName: "P",
				Classes: []string{"active"},
				Attributes: [][2]string{
					{"title", "Open"},
				},
			},
		},
		{
			Input: "p[title='value []'][data-x='escaped \\' quote']",
			Error: false,
			Result: Result{
				TagName: "p",
				Attributes: [][2]string{
					{"title", "value []"},
					{"data-x", "escaped ' quote"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			result, err := Parse(tt.Input)
			if err == nil && !tt.Error && !reflect.DeepEqual(*result, tt.Result) {
				t.Errorf("got %#v; expected %#v", *result, tt.Result)
			} else if err == nil && !tt.Error && reflect.DeepEqual(*result, tt.Result) && result.String() != tt.Input {
				t.Errorf("got String() = %#v; expected %#v", result.String(), tt.Input)
			} else if err == nil && tt.Error {
				t.Errorf("no error, but expected one")
			} else if err != nil && !tt.Error {
				t.Errorf("got error: %s", err)
			}
		})
	}
}

func TestResult_String(t *testing.T) {
	var r *Result
	if str, expected := r.String(), "<nil>"; str != expected {
		t.Fatalf("got %#v; expected %#v", str, expected)
	}
}
