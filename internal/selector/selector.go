package selector

import (
	"errors"
	"io"
	"strings"
)

type Result struct {
	TagName    string
	ID         string
	Classes    []string
	Attributes [][2]string
}

func valueNeedsEscaped(s string) bool {
	return strings.IndexAny(s, "']") != -1
}

func (r *Result) String() string {
	if r == nil {
		return "<nil>"
	}
	var b strings.Builder
	b.WriteString(r.TagName)
	if r.ID != "" {
		b.WriteByte('#')
		b.WriteString(r.ID)
	}
	for _, class := range r.Classes {
		b.WriteByte('.')
		b.WriteString(class)
	}
	for _, attr := range r.Attributes {
		b.WriteByte('[')
		b.WriteString(attr[0])
		if attr[1] != "" {
			b.WriteByte('=')
			if valueNeedsEscaped(attr[1]) {
				b.WriteByte('\'')
				for _, ch := range attr[1] {
					switch ch {
					case '\'':
						b.WriteByte('\\')
						b.WriteRune(ch)
					default:
						b.WriteRune(ch)
					}
				}
				b.WriteByte('\'')
			} else {
				b.WriteString(attr[1])
			}
		}
		b.WriteByte(']')
	}
	return b.String()
}

type Error struct {
	Err      error
	Selector string
}

func newError(underlying error, selector string) *Error {
	return &Error{
		Err:      underlying,
		Selector: selector,
	}
}

func (e *Error) Error() string {
	if e.Selector != "" {
		return e.Err.Error() + " (" + e.Selector + ")"
	}
	return e.Err.Error()
}

var (
	ErrInvalid      = errors.New("invalid selector")
	ErrInvalidID    = errors.New("invalid ID")
	ErrInvalidClass = errors.New("invalid class")
	ErrInvalidAttr  = errors.New("invalid attribute")
)

func Parse(s string) (*Result, error) {
	if s == "" {
		return &Result{}, nil
	}

	var r strings.Reader
	r.Reset(s)

	result := new(Result)
	result.TagName = nextID(&r, "#.[")

	if ch, _, _ := r.ReadRune(); ch == '#' {
		result.ID = nextID(&r, "#.[")
		if result.ID == "" {
			return nil, newError(ErrInvalidID, s)
		}
	} else {
		r.UnreadRune()
	}

	for {
		if ch, _, _ := r.ReadRune(); ch == '.' {
			class := nextID(&r, "#.[")
			if class == "" {
				return nil, newError(ErrInvalidClass, s)
			}
			result.Classes = append(result.Classes, class)
		} else {
			r.UnreadRune()
			break
		}
	}

	for {
		if ch, _, _ := r.ReadRune(); ch == '[' {
			var attrKey, attrValue string
			attrKey = nextID(&r, "=]")
			if attrKey == "" || r.Len() == 0 {
				return nil, newError(ErrInvalidAttr, s)
			}
			if ch, _, _ := r.ReadRune(); ch == ']' {
				// empty attribute value
			} else if ch == '=' {

				if ch, _, _ := r.ReadRune(); ch == '\'' {
					attrValue = nextQuotedValue(&r)
				} else {
					r.UnreadRune()
					attrValue = nextUnquotedValue(&r)
				}

				if ch, _, err := r.ReadRune(); err != nil || ch != ']' {
					return nil, newError(ErrInvalidAttr, s)
				}
			} else {
				return nil, newError(ErrInvalidAttr, s)
			}
			result.Attributes = append(result.Attributes, [2]string{attrKey, attrValue})
		} else {
			r.UnreadRune()
			break
		}
	}

	if r.Len() != 0 {
		return nil, newError(ErrInvalid, s)
	}

	return result, nil
}

func nextID(r *strings.Reader, exceptCharset string) string {
	if len(exceptCharset) == 0 {
		panic("empty exceptCharset")
	}

	var s strings.Builder
	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF || strings.ContainsAny(string(ch), exceptCharset) {
			r.UnreadRune()
			break
		}
		s.WriteRune(ch)
	}
	return s.String()
}

func nextQuotedValue(r *strings.Reader) string {
	var b strings.Builder
	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF {
			return b.String()
		}
		switch ch {
		case '\\':
			if next, _, _ := r.ReadRune(); next == '\'' {
				b.WriteRune(next)
			} else {
				r.UnreadRune()
				return b.String()
			}
		case '\'':
			return b.String()
		default:
			b.WriteRune(ch)
		}
	}
}

func nextUnquotedValue(r *strings.Reader) string {
	var b strings.Builder
	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF {
			break
		}
		if ch == ']' {
			r.UnreadRune()
			break
		}
		b.WriteRune(ch)
	}
	return b.String()
}
