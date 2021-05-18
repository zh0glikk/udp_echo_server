package service

import (
	"context"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3"
	"net"
	"strings"

	"udp_echo_server/internal/config"
)

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

		res := strings.ToUpper(message)

		conn.WriteTo([]byte(res), addr)
	}

	return nil
}


