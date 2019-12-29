package context2

import (
	"fmt"
	"net/http"
)

// Store busca dados
type Store interface {
	Fetch() string
	Cancel()
}

// Server retorna um handler para chamar a Store
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data := make(chan string, 1)

		go func() {
			data <- store.Fetch()
		}()

		select {
		case d := <-data:
			fmt.Fprint(w, d)
		case <-ctx.Done():
			store.Cancel()
		}
	}
}
