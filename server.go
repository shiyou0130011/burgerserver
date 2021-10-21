package burgerserver

import (
	"log"
	"net/http"
	"path"
	"strings"
)

type Filter func(w http.ResponseWriter, r *http.Request, chain FilterChain)

type filterAndPath struct {
	Path   string
	Filter Filter
}

type FilterChain interface {
	Next()
}

type filterChain struct {
	canRunNext bool
}

func (c *filterChain) init() {
	c.canRunNext = false
}

func (c *filterChain) Next() {
	c.canRunNext = true
}

type Server struct {
	serveMux *http.ServeMux
	Filters  []filterAndPath
	Handlers map[string]http.HandlerFunc
	Logger   *log.Logger
}

var defaultServer = DefaultServer()

func NewServer() *Server {
	s := &Server{
		Logger:   log.Default(),
		serveMux: http.DefaultServeMux,
		Handlers: make(map[string]http.HandlerFunc),
	}
	return s
}

func DefaultServer() *Server {
	s := &Server{
		Logger:   log.Default(),
		serveMux: http.DefaultServeMux,
		Handlers: make(map[string]http.HandlerFunc),
	}
	s.AddFilter("/", func(w http.ResponseWriter, r *http.Request, chain FilterChain) {
		s.Logger.Printf(`%s: %s  %v`, r.RemoteAddr, r.Method, r.URL.String())
		chain.Next()
	})

	return s
}

func (s *Server) AddFilter(urlpattern string, filter Filter) {
	s.Filters = append(s.Filters, filterAndPath{urlpattern, filter})
}

func (s *Server) HandleFunc(urlpattern string, handler http.HandlerFunc) {
	if s.Handlers == nil {
		s.Handlers = make(map[string]http.HandlerFunc)
	}
	s.Handlers[urlpattern] = handler
}

func (s *Server) findFilters(urlpattern string) []Filter {
	handlerFilters := []Filter{}
	for _, f := range s.Filters {
		if strings.HasPrefix(urlpattern, f.Path) {
			s.Logger.Printf("Filter of %s matches", f.Path)
			handlerFilters = append(handlerFilters, f.Filter)
		} else if match, err := path.Match(f.Path, urlpattern); err == nil && match {
			s.Logger.Printf("Filter of %s matches", f.Path)
			handlerFilters = append(handlerFilters, f.Filter)
		} else if err != nil {
			s.Logger.Println(err)
		}
	}

	return handlerFilters
}

func (s *Server) generateHandlerFunc(urlpattern string) http.HandlerFunc {
	s.Logger.Printf("Handle Filters and Handler of %s", urlpattern)
	filters := s.findFilters(urlpattern)
	handler := s.Handlers[urlpattern]
	s.Logger.Printf("Find %d filters of %s", len(filters), urlpattern)

	return func(w http.ResponseWriter, r *http.Request) {
		chain := &filterChain{}
		if len(filters) > 0 {
			for _, filter := range filters {
				chain.init()
				filter(w, r, chain)
				if !chain.canRunNext {
					return
				}
			}
			if !chain.canRunNext {
				return
			}

		}
		handler(w, r)
	}

}

// Generate Server to http.Handler
func (s *Server) ToHttpHandler() http.Handler {
	for urlpattern, _ := range s.Handlers {
		s.serveMux.HandleFunc(
			urlpattern,
			s.generateHandlerFunc(urlpattern),
		)
	}

	return s.serveMux
}
