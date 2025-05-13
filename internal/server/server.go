package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/isaialcantara/toyredis/internal/resp"
)

type Server struct {
	Addr string
}

func New(addr string) *Server { return &Server{addr} }

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

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	defer log.Printf("Connection closed. %+v", conn)
	tokenizer := resp.NewFSMTokenizer(conn)
	parser := resp.NewBasicParser(tokenizer)

	for {
		array, err := parser.NextBulkArray()
		if err != nil {
			var perr resp.ProtocolError
			if errors.As(err, &perr) {
				writeError(conn, perr)
			}
			return
		}

		log.Printf("%+v", array)
		if err := writeOk(conn); err != nil {
			log.Printf("failed to write OK: %v", err)
			return
		}
	}
}

func writeOk(conn net.Conn) error {
	ok := resp.SimpleString("OK")
	_, err := conn.Write(ok.ToRESP())
	return err
}

func writeError(conn net.Conn, perr resp.ProtocolError) {
	_, err := conn.Write(perr.ToRESP())
	if err != nil {
		log.Printf("write error failed: %v", err)
	}

	log.Printf("parser error: %v", perr)
}
