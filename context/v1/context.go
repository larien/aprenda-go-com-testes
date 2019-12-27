package context1

import (
	"fmt"
	"net/http"
)

// Store busca dados
type Store interface {
	Fetch() string
}

// Server retorna um handler para chamar a Store
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, store.Fetch())
	}
}
