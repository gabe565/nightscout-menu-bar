package socket

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
)

const (
	timeout = time.Second

	// CommandSettings asks the running instance to open the settings window.
	CommandSettings = "settings"
)

// ErrRunning is returned by Listen when another instance already holds the socket.
var ErrRunning = errors.New("another instance is already running")

// Path returns the per-user socket path.
func Path() string {
	name := "nightscout.sock"
	// Windows has no uid, but its temp directory is per-user
	if uid := os.Getuid(); uid != -1 {
		name = "nightscout-" + strconv.Itoa(uid) + ".sock"
	}

	if runtime.GOOS == "linux" {
		if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
			return filepath.Join(dir, name)
		}
	}
	return filepath.Join(os.TempDir(), name)
}

// New creates a socket that serves the latest reading and doubles as a
// single-instance guard. onCommand is called with commands sent by other instances.
func New(conf *config.Config, onCommand func(string)) *Socket {
	return &Socket{
		config:    conf,
		path:      Path(),
		onCommand: onCommand,
	}
}

type Socket struct {
	config    *config.Config
	path      string
	onCommand func(string)
	listener  *net.UnixListener
	last      *nightscout.Properties
	mu        sync.RWMutex
	wg        sync.WaitGroup
}

func (s *Socket) Format(last *nightscout.Properties) string {
	data := s.config.Data()

	return last.Bgnow.DisplayBg(data.Units) + "," +
		last.Bgnow.Arrow(data.Arrows) + "," +
		last.Delta.Display(data.Units) + "," +
		last.Bgnow.Mills.Relative(data.Advanced.RoundAge) + "," +
		strconv.Itoa(int(time.Since(last.Bgnow.Mills.Time).Seconds())) +
		"\n"
}

func (s *Socket) Write(last *nightscout.Properties) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = last
}

// Listen claims the socket. If another instance is already listening, it is
// asked to open the settings window and ErrRunning is returned.
func (s *Socket) Listen() error {
	if err := s.Close(); err != nil {
		slog.Error("Failed to close socket", "error", err)
	}

	addr := &net.UnixAddr{Name: s.path, Net: "unix"}
	var listener *net.UnixListener
	var err error
	for range 2 {
		if listener, err = net.ListenUnix("unix", addr); err == nil {
			break
		}

		if sendErr := Send(s.path, CommandSettings); sendErr == nil {
			return ErrRunning
		}

		// Nothing answered, so the socket was left behind by an instance that didn't exit cleanly
		slog.Debug("Removing stale socket", "path", s.path)
		if err := os.Remove(s.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(s.path, 0o700); err != nil {
			_ = listener.Close()
			return err
		}
	}

	s.mu.Lock()
	s.listener = listener
	s.mu.Unlock()

	go func() {
		for {
			conn, err := listener.AcceptUnix()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					slog.Error("Failed to accept connection", "error", err)
				}
				return
			}

			s.wg.Go(func() {
				s.handle(conn)
			})
		}
	}()

	return nil
}

// handle writes the latest reading, then reads an optional command.
func (s *Socket) handle(conn *net.UnixConn) {
	defer func() {
		_ = conn.Close()
	}()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	s.mu.RLock()
	last := s.last
	s.mu.RUnlock()

	// An empty line lets readers fail fast instead of waiting for a timeout
	reading := "\n"
	if s.config.Data().Socket.Enabled && last != nil {
		reading = s.Format(last)
	}
	_, _ = io.WriteString(conn, reading)
	_ = conn.CloseWrite()

	line, err := bufio.NewReader(io.LimitReader(conn, 64)).ReadString('\n')
	if err != nil && (!errors.Is(err, io.EOF) || line == "") {
		return
	}
	if s.onCommand != nil {
		if cmd := strings.TrimSpace(line); cmd != "" {
			slog.Debug("Received socket command", "command", cmd)
			s.onCommand(cmd)
		}
	}
}

// Send connects to the socket at path and sends a command.
func Send(path, cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "unix", path)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	if _, err := io.WriteString(conn, cmd+"\n"); err != nil {
		return err
	}
	if c, ok := conn.(*net.UnixConn); ok {
		_ = c.CloseWrite()
	}
	// Wait for the other side to close so the command is known to be handled
	_, err = io.Copy(io.Discard, conn)
	return err
}

func (s *Socket) Close() error {
	s.mu.Lock()
	listener := s.listener
	s.listener = nil
	s.mu.Unlock()

	var err error
	if listener != nil {
		err = listener.Close()
	}
	s.wg.Wait()
	return err
}
