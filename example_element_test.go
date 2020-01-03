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
