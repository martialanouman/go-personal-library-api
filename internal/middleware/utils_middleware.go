package middleware

import (
	"context"
	"fmt"
	"net/http"
)

type Pagination struct {
	Page int
	Take int
}

type UtilsMiddleware struct{}

type PaginationContextKey string

const (
	PaginationKey = PaginationContextKey("pagination")
)

func SetPagination(r *http.Request, p Pagination) *http.Request {
	ctx := context.WithValue(r.Context(), PaginationKey, p)
	return r.WithContext(ctx)
}

func GetPagination(r *http.Request) Pagination {
	p, ok := r.Context().Value(PaginationKey).(Pagination)
	if !ok {
		panic("could not get pagination from request context")
	}

	return p
}

func NewUtilsMiddleware() UtilsMiddleware {
	return UtilsMiddleware{}
}

func (m *UtilsMiddleware) GetPagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page := 1
		take := 10

		if p := q.Get("page"); p != "" {
			fmt.Sscanf(p, "%d", &page)
		}

		if t := q.Get("take"); t != "" {
			fmt.Sscanf(t, "%d", &take)
		}

		pagination := Pagination{
			Page: page,
			Take: take,
		}

		r = SetPagination(r, pagination)
		next.ServeHTTP(w, r)
	})
}
