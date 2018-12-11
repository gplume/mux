package middle

import (
	"net/http"
)

// Wrapper type (it gets its name from the adapter pattern — also known as the decorator pattern)
// above is a function that both takes in and returns an http.Handler. This is the essence of the wrapper;
// we will pass in an existing http.Handler, the Adapter will adapt it, and return a new (probably wrapped) http.Handler
// for us to use in its place. So far this is not much different from just wrapping http.HandlerFunc types,
// however, now, we can instead write functions that themselves return an Wrapper.
type Wrapper func(http.Handler) http.Handler

// Ware is our function takes the handler you want to adapt,
// and a list of our Ware types. The result of our Middlewares
// is an acceptable Ware. Our Ware function will simply iterate over all wrappers,
//  calling them one by one (in reverse order) in a chained manner, returning the result of the first wrapper.
func Ware(h http.Handler, wrappers ...Wrapper) http.Handler {
	// reverse order:
	for _, warp := range wrappers {
		h = warp(h)
	}
	// straight order:
	// for i := len(wrappers) - 1; i >= 0; i-- {
	// 	h = wrappers[i](h)
	// }
	return h
}
