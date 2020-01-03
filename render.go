package m

import (
	"io"
	"strings"
)

// Render writes the HTML of element to w.
//
// A non-nil error is returned if the element could not be successfully written.
func Render(w io.Writer, element Element) error {
	if element == nil {
		return nil
	}
	if internal, ok := element.(internalElement); ok {
		if err := internal.renderHTML(w); err != nil {
			return err
		}
		return nil
	}
	return Render(w, element.Element())
}

// RenderString returns the HTML of element.
func RenderString(element Element) string {
	var b strings.Builder
	Render(&b, element)
	return b.String()
}
