package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type HandlerSpec struct {
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]HandlerSpec
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]HandlerSpec),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method string, path string, handler http.HandlerFunc) {
	s.Handlers[method+" "+path] = HandlerSpec{Method: method, Path: path, HandlerFunc: handler}
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)

	for _, httpHandler := range s.Handlers {
		s.Router.Method(httpHandler.Method, httpHandler.Path, httpHandler.HandlerFunc)
	}

	http.ListenAndServe(":"+s.WebServerPort, s.Router)
}
