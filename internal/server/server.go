package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/isaialcantara/toyredis/internal/command"
	"github.com/isaialcantara/toyredis/internal/resp"
)

type Server struct {
	dispatcher command.Dispatcher
	Addr       string
}

func New(addr string, dispatcher command.Dispatcher) *Server {
	return &Server{dispatcher, addr}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}
	log.Printf("Started server on: %s", s.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer log.Println("Connection closed")
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection. %v", err)
		}
	}()

	parser := resp.NewRESPParser(conn)

	for {
		bulkArray, err := parser.NextBulkArray()
		if err != nil {
			var respErr resp.RESPError
			if errors.As(err, &respErr) {
				writeError(conn, respErr)
			}
			return
		}

		response := s.dispatcher.Dispatch(bulkArray)
		if err := writeResponse(conn, response); err != nil {
			log.Printf("failed to write response: %v", err)
			return
		}
	}
}

func writeResponse(conn net.Conn, response []byte) error {
	_, err := conn.Write(response)
	return err
}

func writeError(conn net.Conn, respErr resp.RESPType) {
	_, err := conn.Write(respErr.ToRESP())
	if err != nil {
		log.Printf("write error failed: %v", err)
	}

	log.Printf("parser error: %v", respErr)
}
