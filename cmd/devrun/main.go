package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

type serviceCommand struct {
	name    string
	args    []string
	dir     string
	waitFor string
}

func main() {
	moduleRoot, err := findModuleRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to locate go module root: %v\n", err)
		os.Exit(1)
	}

	services := []serviceCommand{
		{name: "user-rpc", args: []string{"run", "."}, dir: "apps/user/rpc"},
		{name: "user-api", args: []string{"run", "."}, dir: "apps/user/api"},
		{name: "im-ws", args: []string{"run", "."}, dir: "apps/im/ws"},
		{name: "task-mq", args: []string{"run", "."}, dir: "apps/task/mq", waitFor: "127.0.0.1:10090"},
		{name: "social-rpc", args: []string{"run", "."}, dir: "apps/social/rpc"},
		{name: "social-api", args: []string{"run", "."}, dir: "apps/social/api"},
		{name: "im-rpc", args: []string{"run", "."}, dir: "apps/im/rpc"},
		{name: "im-api", args: []string{"run", "."}, dir: "apps/im/api"},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	errs := make(chan error, len(services))

	for _, svc := range services {
		wg.Add(1)
		go func(s serviceCommand) {
			defer wg.Done()
			errs <- runService(ctx, cancel, s, moduleRoot)
		}(svc)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	exitCode := 0
	for err := range errs {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

func runService(ctx context.Context, cancel context.CancelFunc, svc serviceCommand, moduleRoot string) error {
	if svc.waitFor != "" {
		if err := waitForPort(ctx, svc.waitFor, 30*time.Second); err != nil {
			cancel()
			return fmt.Errorf("%s dependency %s not ready: %w", svc.name, svc.waitFor, err)
		}
	}
	cmd := exec.CommandContext(ctx, "go", svc.args...)
	cmd.Stdout = newPrefixWriter(os.Stdout, svc.name)
	cmd.Stderr = newPrefixWriter(os.Stderr, svc.name)
	cmd.Env = os.Environ()
	workDir := moduleRoot
	if svc.dir != "" {
		workDir = filepath.Join(moduleRoot, svc.dir)
	}
	cmd.Dir = workDir

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start %s: %w", svc.name, err)
	}

	if err := cmd.Wait(); err != nil {
		if ctx.Err() != nil {
			return nil
		}
		cancel()
		return fmt.Errorf("%s exited: %w", svc.name, err)
	}

	cancel()
	return fmt.Errorf("%s exited unexpectedly", svc.name)
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found starting from %s", dir)
		}
		dir = parent
	}
}

// prefixWriter tags each log line with the service name so interleaved output stays readable.
type prefixWriter struct {
	dst     io.Writer
	prefix  []byte
	mu      sync.Mutex
	newLine bool
}

func newPrefixWriter(dst io.Writer, name string) io.Writer {
	return &prefixWriter{
		dst:     dst,
		prefix:  []byte("[" + name + "] "),
		newLine: true,
	}
}

func (w *prefixWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	written := 0
	for len(p) > 0 {
		if w.newLine {
			if _, err := w.dst.Write(w.prefix); err != nil {
				return written, err
			}
			w.newLine = false
		}

		idx := bytes.IndexByte(p, '\n')
		if idx == -1 {
			n, err := w.dst.Write(p)
			written += n
			return written, err
		}

		n, err := w.dst.Write(p[:idx+1])
		written += n
		if err != nil {
			return written, err
		}

		p = p[idx+1:]
		w.newLine = true
	}

	return written, nil
}

func waitForPort(ctx context.Context, addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		dialer := net.Dialer{Timeout: time.Second}
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for %s: %w", addr, err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}
