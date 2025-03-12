package disignpattern

import "fmt"

type Server struct {
	host string
	port int
}

type Option func(*Server)

func NewServer(opts ...Option) *Server {
	s := &Server{host: "localhost", port: 8080}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func OptionExample() {
	server := NewServer(WithHost("0.0.00."), WithPort(9999))
	fmt.Println(server)
}
