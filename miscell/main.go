package mai

import (
	"net/http"
	"regexp"
)

/*
	closure in https://go.dev/doc/articles/wiki/
	take a function as param return another function this funciton is a function literals
	this literals function is same type of return function but make a side effect then call the
	passed in function

*/

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

/*
the purpose of makeHandler is to reduce deplication in editHandler and viewHandler they both run validation before
renderTemplate.
*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
}

/*
	transform differet functions to an interface

	if bounch of functions doing same thing as interface implies

	you can trasnsfrom these functions into interface

	see color package  in go std lib
*/

type Model interface {
	Convert(c Color) Color
}

// ModelFunc returns a Model that invokes f to implement the conversion.
func ModelFunc(f func(Color) Color) Model {
	// 1.return a struct which implements the Model interface
	return &modelFunc{f}
}

// 2. struct wraps the function
type modelFunc struct {
	f func(Color) Color
}

// 3. the struct implement the interface by calling the wrapped function
func (m *modelFunc) Convert(c Color) Color {
	return m.f(c)
}

var (
	RGBAModel   Model = ModelFunc(rgbaModel)
	RGBA64Model Model = ModelFunc(rgba64Model)
)

func rgbaModel(c Color) Color {
	if _, ok := c.(RGBA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func rgba64Model(c Color) Color {
	if _, ok := c.(RGBA64); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}
