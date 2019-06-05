package web

import (
	"fmt"
	"net/http"
	"strings"
)

type GetHandler interface {
	Get(http.ResponseWriter, *http.Request)
}

type HeadHandler interface {
	Head(http.ResponseWriter, *http.Request)
}

type PostHandler interface {
	Post(http.ResponseWriter, *http.Request)
}

type PutHandler interface {
	Put(http.ResponseWriter, *http.Request)
}

type DeleteHandler interface {
	Delete(http.ResponseWriter, *http.Request)
}

// MethodHandler
type methodHandler struct {
	methods []string
	handler []http.HandlerFunc
}

// MethodHandler as generic error
func (m methodHandler) Error() string {
	return fmt.Sprintf("Allow: %s", strings.Join(m.methods, ", "))
}

// MethodHandler as MethodNotAllowed error
func (m methodHandler) Methods() []string {
	return m.methods
}

// MethodHandler as http.Handler
func (m methodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		for i, k := range m.methods {
			if k == r.Method {
				m.handler[i](w, r)
				return
			}
		}
	}

	panic(m) // 405
}

// MethodHandler constructor
func (m methodHandler) addMethodHandlerFunc(method string, h http.HandlerFunc) {
	m.methods = append(m.methods, method)
	m.handler = append(m.handler, h)
}

func MethodHandler(h interface{}) http.Handler {
	var m methodHandler
	var get GetHandler

	// GET
	if o, ok := h.(GetHandler); ok {
		get = o
		m.addMethodHandlerFunc("GET", o.Get)
	}
	// HEAD
	if o, ok := h.(HeadHandler); ok {
		m.addMethodHandlerFunc("HEAD", o.Head)
	} else if get != nil {
		m.addMethodHandlerFunc("HEAD", get.Get)
	}
	// POST
	if o, ok := h.(PostHandler); ok {
		m.addMethodHandlerFunc("POST", o.Post)
	}
	// PUT
	if o, ok := h.(PutHandler); ok {
		m.addMethodHandlerFunc("PUT", o.Put)
	}
	// DELETE
	if o, ok := h.(DeleteHandler); ok {
		m.addMethodHandlerFunc("DELETE", o.Delete)
	}

	return m
}
