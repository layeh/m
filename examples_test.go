package m_test

import (
	"fmt"
	"os"
	"strings"

	. "layeh.com/m"
)

func ExampleRender() {
	el := M("p",
		T("Hello World"),
	)
	if err := Render(os.Stdout, el); err != nil {
		panic(err)
	}
	// Output:
	// <p>Hello World</p>
}

func ExampleRenderString() {
	el := M("p",
		T("Hello World"),
	)
	fmt.Println(RenderString(el))
	// Output:
	// <p>Hello World</p>
}

func ExampleAttr() {
	el := M("p", Attr("data-id", "ff34"),
		T("User"),
	)
	fmt.Println(RenderString(el))
	// Output:
	// <p data-id="ff34">User</p>
}

func ExampleAttrf() {
	el := M("div", Attrf("id", "user-%d", 10),
		T("User section"),
	)
	fmt.Println(RenderString(el))
	// Output:
	// <div id="user-10">User section</div>
}

func ExampleDocument() {
	doc := Document(
		M("head",
			M("title", T("Hello World")),
		),
		M("body"),
	)
	if err := Render(os.Stdout, doc); err != nil {
		panic(err)
	}
	// Output:
	// <!DOCTYPE html>
	// <head><title>Hello World</title></head><body></body>
}

func ExampleF() {
	el := M("p",
		F("Hello, %s", "World"),
	)
	fmt.Println(RenderString(el))
	// Output:
	// <p>Hello, World</p>
}

func ExampleFor() {
	el := M("ul",
		For(0, 3, 1, func(i int) Element {
			return M("li", F("%d", i*10))
		}),
	)
	fmt.Println(RenderString(el))
	// Output:
	// <ul><li>0</li><li>10</li><li>20</li></ul>
}

func ExampleGroup() {
	users := []string{
		"Alice",
		"Bob",
		"Bill",
		"Eve",
	}

	el := Group(len(users), func(i, j int) bool {
		// Group users by first letter of name
		return users[i][0] == users[j][0]
	}, func(i, j int) Element {
		return M("p",
			T(strings.Join(users[i:j], ", ")),
		)
	})
	fmt.Println(RenderString(el))
	// Output:
	// <p>Alice</p><p>Bob, Bill</p><p>Eve</p>
}

func ExampleIf() {
	el := Range(5, func(i int) Element {
		return M("p",
			F("%d", i),
			If(i%2 == 0, T("!")),
		)
	})
	fmt.Println(RenderString(el))
	// Output:
	// <p>0!</p><p>1</p><p>2!</p><p>3</p><p>4!</p>
}

func ExampleIfElse() {
	el := Range(5, func(i int) Element {
		return M("p",
			F("%d", i),
			IfElse(i%2 == 0, T("!"), T("?")),
		)
	})
	fmt.Println(RenderString(el))
	// Output:
	// <p>0!</p><p>1?</p><p>2!</p><p>3?</p><p>4!</p>
}

func ExampleM() {
	el := M("h1#headline.active.etc[data-id=3]",
		T("Hello World"),
	)
	fmt.Println(RenderString(el))
	// Output:
	// <h1 id="headline" class="active etc" data-id="3">Hello World</h1>
}
