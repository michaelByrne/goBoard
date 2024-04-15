package session

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionKey string

const (
	key SessionKey = "_session_store"
)

func Get(name string, r *http.Request) (*sessions.Session, error) {
	s := r.Context().Value(key)
	if s == nil {
		return nil, fmt.Errorf("%q session store not found", key)
	}

	store := s.(sessions.Store)
	session, err := store.Get(r, name)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func SessionMiddleware(store sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, key, store)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
