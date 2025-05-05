package socket

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"strconv"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
	"gabe565.com/nightscout-menu-bar/internal/util"
)

func New(conf *config.Config) *Socket {
	l := &Socket{
		config: conf,
	}
	l.reloadConfig()

	conf.AddCallback(l.reloadConfig)
	return l
}

type Socket struct {
	config   *config.Config
	listener *net.UnixListener
	last     *nightscout.Properties
}

func (s *Socket) Format(last *nightscout.Properties) string {
	switch s.config.Socket.Format {
	case config.SocketFormatCSV:
		return last.Bgnow.DisplayBg(s.config.Units) + "," +
			last.Bgnow.Arrow(s.config.Arrows) + "," +
			last.Delta.Display(s.config.Units) + "," +
			last.Bgnow.Mills.Relative(s.config.Advanced.RoundAge) + "," +
			strconv.Itoa(int(time.Since(last.Bgnow.Mills.Time).Seconds())) +
			"\n"
	default:
		slog.Error("Unknown socket format", "value", s.config.Socket.Format)
		return ""
	}
}

func (s *Socket) reloadConfig() {
	if err := s.Close(); err != nil {
		slog.Error("Failed to close socket", "error", err)
	}

	if s.config.Socket.Enabled {
		if err := s.Listen(); err != nil {
			slog.Error("Failed to listen on socket", "error", err)
		}
	}
}

func (s *Socket) Write(last *nightscout.Properties) {
	s.last = last
}

func (s *Socket) Listen() error {
	if err := s.Close(); err != nil {
		slog.Error("Failed to close socket", "error", err)
	}

	var err error
	s.listener, err = net.ListenUnix("unix", &net.UnixAddr{
		Name: util.ResolvePath(s.config.Socket.Path),
		Net:  "unix",
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					slog.Error("Failed to accept connection", "error", err)
				}
				return
			}

			go func() {
				defer func() {
					_ = conn.Close()
				}()
				_ = conn.SetDeadline(time.Now().Add(time.Second))
				_, _ = io.WriteString(conn, s.Format(s.last))
			}()
		}
	}()

	return nil
}

func (s *Socket) Close() error {
	var err error
	if s.listener != nil {
		err = s.listener.Close()
	}
	s.listener = nil
	return err
}
