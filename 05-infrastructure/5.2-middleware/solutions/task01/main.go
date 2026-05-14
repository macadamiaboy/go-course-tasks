package main

import (
	"fmt"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func testMiddleware(name string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[%s] before\n", name)
			next.ServeHTTP(w, r)
			fmt.Printf("[%s] after\n", name)
		})
	}
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "hello, world")
}

func main() {
	mux := http.NewServeMux()

	handler := Chain(
		http.HandlerFunc(helloHandler),
		testMiddleware("mw1"),
		testMiddleware("mw2"),
	)
	mux.Handle("GET /hello", handler)

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
