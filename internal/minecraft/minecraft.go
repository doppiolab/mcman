package minecraft

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source minecraft.go -destination minecraft_mock.go -package "minecraft" MinecraftServer

type MinecraftServer interface {
	// Start subprocess and returns stdout and stderr channel.
	Start() (chan string, chan string, error)

	// If server is already exited or stop command is sent successfully, return nil, otherwise return error.
	Stop() error

	// Send command to the minecraft server process.
	PutCommand(cmd string) error

	// Get minecraft server process
	GetProcess() *exec.Cmd
}

type minecraftServer struct {
	cfg *config.MinecraftConfig

	svrProcess *exec.Cmd
	stdinPipe  io.WriteCloser
	stdoutPipe io.ReadCloser
	stderrPipe io.ReadCloser

	stdinLock sync.Mutex
}

func NewMinecraftServer(cfg *config.MinecraftConfig) (MinecraftServer, error) {
	minecraftCommand := createMinecraftCommandArgs(cfg)
	log.Info().Msgf("launch minecraft server. cmd: %s %s", cfg.JavaCommand, strings.Join(minecraftCommand, " "))
	mcCmd := exec.Command(cfg.JavaCommand, minecraftCommand...)
	// prevent shell to send signal to subprocess
	mcCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdinPipe, err := mcCmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get stdin pipe")
	}
	stdoutPipe, err := mcCmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get stdout pipe")
	}
	stderrPipe, err := mcCmd.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get stderr pipe")
	}

	if cfg.WorkingDir != "" {
		mcCmd.Dir = cfg.WorkingDir
		if err := maybeCreateWorkingDir(mcCmd.Dir); err != nil {
			return nil, err
		}
	}

	return &minecraftServer{
		svrProcess: mcCmd,
		cfg:        cfg,
		stdinPipe:  stdinPipe,
		stdoutPipe: stdoutPipe,
		stderrPipe: stderrPipe,
	}, nil
}

func (s *minecraftServer) Start() (chan string, chan string, error) {
	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	if s.cfg.SkipStartForDebug {
		log.Warn().Msg("Skip start minecraft server for debug purpose. Do not use this flag for the real server.")
		return stdoutChan, stderrChan, nil
	}

	if err := s.svrProcess.Start(); err != nil {
		return nil, nil, errors.Wrap(err, "cannot start server subprocess")
	}

	go streamReaderToChan(s.stdoutPipe, stdoutChan)
	go streamReaderToChan(s.stderrPipe, stderrChan)

	return stdoutChan, stderrChan, nil
}

func (s *minecraftServer) Stop() error {
	if s.cfg.SkipStartForDebug {
		log.Warn().Msg("Skip stop minecraft server for debug purpose.")
		return nil
	}

	if err := s.PutCommand("stop"); err != nil {
		return errors.Wrap(err, "cannot stop server")
	}

	if err := s.svrProcess.Wait(); err != nil {
		return errors.Wrap(err, "cannot wait server process")
	}

	return nil
}

func (s *minecraftServer) PutCommand(cmd string) error {
	s.stdinLock.Lock()
	defer s.stdinLock.Unlock()

	cmd = strings.Trim(cmd, " \n\r")
	n, err := io.WriteString(s.stdinPipe, cmd+"\n")
	log.Info().Int("n", n).Msgf("sent command \"%s\"", cmd)
	if err != nil {
		return errors.Wrap(err, "cannot put command")
	}
	return nil
}

func (s *minecraftServer) GetProcess() *exec.Cmd {
	return s.svrProcess
}

// Create working directory if it doesn't exist.
func maybeCreateWorkingDir(path string) error {
	fileInfo, err := os.Stat(path)

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "cannot get working dir")
		}

		log.Info().Str("path", path).Msg("create working dir for server")
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Wrap(err, "cannot create working dir")
		}
	} else if !fileInfo.IsDir() {
		return errors.New("working dir does exist but not a directory")
	}

	return nil
}

// Create minecraft server launch command from minecraft config.
func createMinecraftCommandArgs(cfg *config.MinecraftConfig) []string {
	result := []string{}
	result = append(result, cfg.JavaOptions...)
	result = append(result, "-jar")
	result = append(result, cfg.JarPath)
	result = append(result, cfg.Args...)
	return result
}

// stream io.Reader to write only chan.
//
// If this function met io.EOF.
func streamReaderToChan(in io.Reader, out chan<- string) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	close(out)
}
