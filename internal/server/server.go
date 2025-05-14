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
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection. %v", err)
		}
	}()

	defer log.Printf("Connection closed. %+v", conn)
	tokenizer := resp.NewFSMTokenizer(conn)
	parser := resp.NewBasicParser(tokenizer)

	for {
		bulkArray, err := parser.NextBulkArray()
		if err != nil {
			var respErr resp.RESPType
			if errors.As(err, &respErr) {
				writeError(conn, respErr)
			}
			return
		}

		response, err := command.DispatchCommand(bulkArray)
		if err != nil {
			var respErr resp.RESPType
			if errors.As(err, &respErr) {
				writeError(conn, respErr)
			}
		} else {
			if err := writeResponse(conn, response); err != nil {
				log.Printf("failed to write response: %v", err)
				return
			}
		}
	}
}

func writeResponse(conn net.Conn, response resp.RESPType) error {
	_, err := conn.Write(response.ToRESP())
	return err
}

func writeError(conn net.Conn, respErr resp.RESPType) {
	_, err := conn.Write(respErr.ToRESP())
	if err != nil {
		log.Printf("write error failed: %v", err)
	}

	log.Printf("parser error: %v", respErr)
}
