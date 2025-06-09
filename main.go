package main

import (
	"net/http"

	"github.com/romsar/gonertia/v2"
)

func main() {
	inertia, err := NewInertia()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/home", inertia.Middleware(HomeHandler(inertia)))

	mux.Handle("/public/build/assets/", 
    http.StripPrefix("/public/build/assets/", 
        http.FileServer(http.Dir("./public/build/assets")),
    ),
)

	http.ListenAndServe(":8080", mux)
}

func HomeHandler(i *gonertia.Inertia) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		i.Render(w, r, "index", gonertia.Props{
            "some": "data",
        })
	}

	return http.HandlerFunc(f)
}
