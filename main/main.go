package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/akamensky/argparse"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func sshMain(host *string, port *int) {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(*host, strconv.Itoa(*port))),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", *host, "port", *port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// pty, _, _ := s.Pty()

	model := initialModel()

	return model, []tea.ProgramOption{tea.WithAltScreen()}
}

func standaloneMain() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main() {
	parser := argparse.NewParser("ssh-ttt", "Play Tic-tac-toe in your terminal!")
	standalone := parser.Flag("s", "standalone", &argparse.Options{Default: false, Help: "Starts the game in current terminal instead of starting the SSH server"})
	host := parser.String("l", "listen", &argparse.Options{Default: "0.0.0.0", Help: "Interface to listen on"})
	port := parser.Int("p", "port", &argparse.Options{Default: 23234, Help: "Port to listen on"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if *standalone {
		standaloneMain()
	} else {
		sshMain(host, port)
	}
}
