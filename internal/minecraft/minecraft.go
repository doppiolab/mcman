package minecraft

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source minecraft.go -destination minecraft_mock.go -package "minecraft" MinecraftServer

type MinecraftServer interface {
	// Start subprocess and stream stdout and stderr to channels.
	Start() (stdout chan string, stderr chan string, err error)

	// If server is already exited or stop command is sent successfully, return nil, otherwise return error.
	Stop() error

	// Send command to the minecraft server process.
	PutCommand(cmd string) error
}

type minecraftServer struct {
	cfg *config.MinecraftConfig

	svrProcess *exec.Cmd
	stdinPipe  io.WriteCloser
	stdoutPipe io.ReadCloser
	stderrPipe io.ReadCloser
}

func NewMinecraftServer(cfg *config.MinecraftConfig) (MinecraftServer, error) {
	minecraftCommandString := createMinecraftCommand(cfg)
	log.Info().Msgf("launch minecraft server. cmd: %s", minecraftCommandString)
	mcCmd := exec.Command("bash", "-c", minecraftCommandString)

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

func (s *minecraftServer) Start() (stdout chan string, stderr chan string, err error) {
	if err := s.svrProcess.Start(); err != nil {
		return nil, nil, errors.Wrap(err, "cannot start server subprocess")
	}

	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	go streamReaderToChan(s.stdoutPipe, stdoutChan)
	go streamReaderToChan(s.stderrPipe, stderrChan)

	return stdoutChan, stderrChan, nil
}

func (s *minecraftServer) Stop() error {
	if err := s.PutCommand("stop"); err != nil {
		return errors.Wrap(err, "cannot stop server")
	}

	if err := s.svrProcess.Wait(); err != nil {
		return errors.Wrap(err, "cannot wait server process")
	}

	// clean up all resources
	s.stdinPipe.Close()
	s.stdoutPipe.Close()
	s.stderrPipe.Close()

	return nil
}

func (s *minecraftServer) PutCommand(cmd string) error {
	// TODO(hayeon): add lock to prevent race cond
	cmd = strings.Trim(cmd, " \n\r")
	_, err := s.stdinPipe.Write([]byte(cmd + "\n"))
	if err != nil {
		return errors.Wrap(err, "cannot put command")
	}
	return nil
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
func createMinecraftCommand(cfg *config.MinecraftConfig) string {
	javaOptions := strings.Join(cfg.JavaOptions, " ")
	args := strings.Join(cfg.Args, " ")

	return fmt.Sprintf("%s %s -jar %s %s", cfg.JavaCommand, javaOptions, cfg.JarPath, args)
}

// stream io.Reader to write only chan.
//
// If this function met io.EOF, out channel will be closed.
func streamReaderToChan(in io.Reader, out chan<- string) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	close(out)
}
