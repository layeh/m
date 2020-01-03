package m

import (
	"fmt"
	"html"
	"io"
	"strings"
	"text/template"

	selectorpkg "layeh.com/m/internal/selector"
)

// Element is an object that can be rendered as HTML.
type Element interface {
	// Element returns the Element to be rendered.
	//
	// In most cases, it will return an element constructed with one of
	// the functions in this package.
	Element() Element
}

type internalElement interface {
	renderHTML(w io.Writer) error
}

type htmlElement struct {
	TagName  string
	Elements []Element
	Empty    bool
}

// M returns an element that is an HTML tag that is specified by selector.
//
// Selectors are defined using the following syntax:
//  tagname#id.class-1.class-2[attr-key-1=value][attr-key-2=value]
//
// The element ID (#), class names (.), and attributes ([]) are optional. Multiple
// class names and attributes can be defined, but only one element ID.
// If tagname is not defined, div is used.
//
// The selector should be a constant value. If dynamic values are required for an ID,
// class name, or attribute, omit the dynamic value from the selector string and use
// Attr or Attrf instead.
//
// elements are the children of the HTML tag. If Attr and Attrf values are used,
// they must be the first values included in elements.
//
// The function panics on an invalid selector.
func M(selector string, elements ...Element) Element {
	sel, err := selectorpkg.Parse(selector)
	if err != nil {
		panic(err)
	}

	if sel.TagName == "" {
		sel.TagName = "div"
	} else {
		sel.TagName = strings.ToLower(sel.TagName)
	}

	allElements := make([]Element, 0, 3+len(elements))

	if sel.ID != "" {
		allElements = append(allElements, Attr("id", sel.ID))
	}

	if len(sel.Classes) > 0 {
		allElements = append(allElements, Attr("class", strings.Join(sel.Classes, " ")))
	}

	for _, attr := range sel.Attributes {
		allElements = append(allElements, Attr(attr[0], attr[1]))
	}

	allElements = append(allElements, elements...)

	return &htmlElement{
		TagName:  sel.TagName,
		Elements: allElements,
		Empty:    voidElements[sel.TagName],
	}
}

var voidElements = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}

func (*htmlElement) Element() Element { return nil }

func (e *htmlElement) renderHTML(w io.Writer) error {
	if _, err := io.WriteString(w, "<"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, e.TagName); err != nil {
		return err
	}

	var i int
	// Attributes
	for ; i < len(e.Elements); i++ {
		if attr, ok := e.Elements[i].(*attr); ok {
			if _, err := io.WriteString(w, " "); err != nil {
				return err
			}
			if _, err := io.WriteString(w, attr.Key); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "=\""); err != nil {
				return err
			}
			if _, err := io.WriteString(w, template.HTMLEscapeString(attr.Value)); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "\""); err != nil {
				return err
			}
		} else {
			break
		}
	}

	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}

	if !e.Empty {
		// Children
		for ; i < len(e.Elements); i++ {
			if err := Render(w, e.Elements[i]); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, "</"); err != nil {
			return err
		}
		if _, err := io.WriteString(w, e.TagName); err != nil {
			return err
		}
		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}
	}

	return nil
}

// Document returns an element that renders the HTML5 doctype before elements.
func Document(elements ...Element) Element {
	newSlice := make([]Element, 1+len(elements))
	newSlice[0] = Raw("<!DOCTYPE html>\n")
	copy(newSlice[1:], elements)
	return &slice{
		Elements: newSlice,
	}
}

// S returns an element where each elements are concatenated together.
func S(elements ...Element) Element {
	if len(elements) == 0 {
		return nil
	}

	s := make([]Element, len(elements))
	copy(s, elements)
	return &slice{
		Elements: s,
	}
}

type slice struct {
	Elements []Element
}

func (e *slice) Element() Element { return nil }

func (e *slice) renderHTML(w io.Writer) error {
	for _, element := range e.Elements {
		if err := Render(w, element); err != nil {
			return err
		}
	}
	return nil
}

// Attr returns an HTML element attribute with the given key and value.
//
// The return value is only valid when used as the first elements in calling M.
func Attr(key, value string) Element {
	return &attr{
		Key:   key,
		Value: value,
	}
}

// Attrf returns an HTML element attribute with the given key and value.
// The value is formatted using fmt.
//
// The return value is only valid when used as the first elements in calling M.
func Attrf(key, valueFormat string, x ...interface{}) Element {
	return Attr(key, fmt.Sprintf(valueFormat, x...))
}

type attr struct {
	Key, Value string
}

func (*attr) Element() Element  { return nil }
func (*attr) renderHTML() error { return nil }

// T returns an escaped text element.
func T(text string) Element {
	return &raw{
		Raw: html.EscapeString(text),
	}
}

// Raw returns an element that renders the given HTML unescaped.
func Raw(html string) Element {
	return &raw{
		Raw: html,
	}
}

// F returns an escaped text element that is formatted using fmt.
func F(format string, x ...interface{}) Element {
	return T(fmt.Sprintf(format, x...))
}

type raw struct {
	Raw string
}

func (*raw) Element() Element { return nil }

func (e *raw) renderHTML(w io.Writer) error {
	if _, err := io.WriteString(w, e.Raw); err != nil {
		return err
	}
	return nil
}

// If returns ifTrue if cond is true, nil otherwise.
func If(cond bool, ifTrue Element) Element {
	return IfElse(cond, ifTrue, nil)
}

// IfElse returns ifTrue or ifFalse, depending on the value of cond.
func IfElse(cond bool, ifTrue, ifFalse Element) Element {
	if cond {
		return ifTrue
	}
	return ifFalse
}

// Lazy returns an element that calls fn when it is being rendered.
func Lazy(fn func() Element) Element {
	return &lazyNode{
		Fn: fn,
	}
}

type lazyNode struct {
	Fn func() Element
}

func (*lazyNode) Element() Element { return nil }

func (e *lazyNode) renderHTML(w io.Writer) error {
	if err := Render(w, e.Fn()); err != nil {
		return err
	}
	return nil
}

// Range returns an element that is called for every index from 0 to n.
func Range(n int, fn func(i int) Element) Element {
	return For(0, n, 1, fn)
}

// For returns an element that is called for every index of the loop with the given conditions.
func For(start, end, step int, fn func(i int) Element) Element {
	return &forLoop{
		Start: start,
		End:   end,
		Step:  step,
		Func:  fn,
	}
}

type forLoop struct {
	Start, End, Step int
	Func             func(int) Element
}

func (*forLoop) Element() Element { return nil }

func (e *forLoop) renderHTML(w io.Writer) error {
	if e.Step >= 0 {
		for i := e.Start; i < e.End; i += e.Step {
			if err := Render(w, e.Func(i)); err != nil {
				return err
			}
		}
	} else {
		for i := e.Start; i >= e.End; i += e.Step {
			if err := Render(w, e.Func(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

// Group returns an element that renders contiguous values that match the group function.
//
// From 0 to N (i), the group function is evaluated
//  group(i-1, i)
// If it returns false, the render function is called with the lowest index not yet rendered
// and the current index. If it returns true, the current index is incremented. Render is called
// at the end if any indices remain not yet rendered.
func Group(n int, group func(i, j int) bool, render func(i, j int) Element) Element {
	return &groupEl{
		N:      n,
		Group:  group,
		Render: render,
	}
}

type groupEl struct {
	N      int
	Group  func(i, j int) bool
	Render func(i, j int) Element
}

func (*groupEl) Element() Element { return nil }

func (e *groupEl) renderHTML(w io.Writer) error {
	if e.N == 0 {
		return nil
	}

	lower := 0
	for i := 1; i < e.N; i++ {
		if !e.Group(lower, i) {
			if err := Render(w, e.Render(lower, i)); err != nil {
				return err
			}
			lower = i
		}
	}
	return Render(w, e.Render(lower, e.N))
}
