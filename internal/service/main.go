package service

import (
	"bytes"
	"context"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3"
	"net"
	"os/exec"
	"strings"

	"udp_echo_server/internal/config"
)

const ShellToUse = "/bin/sh"

type service struct {
	logger           *logan.Entry
	cfg 			 config.Config
}

func NewService(cfg config.Config) *service {
	return &service{
		logger:    cfg.Log(),
		cfg:	   cfg,
	}
}

func (s *service) Run(ctx context.Context) error{
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", s.cfg.Listener().Port))
	if err != nil {
		s.logger.WithError(err).Error("failed to resolve udp addr")
		return err
	}

	fmt.Println(fmt.Sprintf("listening on :%v", s.cfg.Listener().Port))

	conn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		s.logger.WithError(err).Error("failed on listening")
		return err
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			s.logger.WithError(err).Error("Failed to read")
		}

		message := string(buf[0:n])

		s.logger.Info(fmt.Sprintf("Recived message %v", message))

		//res := strings.ToUpper(message)
		var res string

		if string(message[0]) == "/"{
			err, stdout, stderr := Shellout(message[1:])
			if err != nil {
				s.logger.Error("failed to exec cmd")
			}

			s.logger.Info(stdout)

			res = stdout + "\n" + stderr
		} else {
			res = strings.ToUpper(message) + "\n"
		}

		conn.WriteTo([]byte(res), addr)
	}

	return nil
}

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

