package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

// https://github.com/bitnami/wait-for-port/blob/8817a80d38e9b0f7b63d15693d201fecce30a4fe/cmd.go
func waitFor(ctx context.Context, port int, timeout time.Duration) error {
	t := time.NewTimer(timeout)
	defer t.Stop()

	check := time.NewTicker(100 * time.Millisecond)
	defer check.Stop()

	addr := fmt.Sprintf("localhost:%d", port)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			return errors.New("timeout waiting for emulator to start")
		case <-check.C:
			if connectable(ctx, port) {
				log.Printf("%s is occupied", addr)
				return nil
			}
		}
	}
}

func connectable(ctx context.Context, port int) bool {
	d := net.Dialer{Timeout: 1 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("localhost:%d", port))
	if err == nil {
		conn.Close()
		return true
	}
	return false
}
