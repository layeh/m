# m [![GoDoc](https://godoc.org/layeh.com/m?status.svg)](https://godoc.org/layeh.com/m)

Package m is an HTML element builder and renderer, inspired by [Mithril.js](https://mithril.js.org/).

```go
package m_test

import (
	"fmt"

	. "layeh.com/m"
)

type Alert struct {
	child Element
}

func NewAlert(child Element) *Alert {
	return &Alert{
		child: child,
	}
}

func (e *Alert) Element() Element {
	return M("div.alert.alert-primary[role=alert]",
		e.child,
	)
}

func ExampleElement() {
	el := NewAlert(T("Access Denied"))
	fmt.Println(RenderString(el))
	// Output:
	// <div class="alert alert-primary" role="alert">Access Denied</div>
}

```

## License

Public domain

## Author

Tim Cooper (<tim.cooper@layeh.com>)
