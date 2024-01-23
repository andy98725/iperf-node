package main

import (
	"bytes"
	"errors"
	"os/exec"
)

func (s *server) runIperfServer() error {
	if s.process != nil {
		return errors.New("iPerf is already running")
	}

	s.process = exec.Command("iperf", "-s -p "+s.config.iperfPort)
	go func() {
		buf := new(bytes.Buffer)
		s.process.Stdout = buf
		s.process.Run()

		results := buf.String()
		s.completeTest(&results)
	}()
	return nil
}

func (s *server) runIperfClient(addr, port string) error {
	if s.process != nil {
		return errors.New("iPerf is already running")
	}

	s.process = exec.Command("iperf", "-c "+addr+" -p "+s.config.iperfPort)
	go func() {
		buf := new(bytes.Buffer)
		s.process.Stdout = buf
		s.process.Run()

		results := buf.String()
		s.completeTest(&results)
	}()
	return nil
}
