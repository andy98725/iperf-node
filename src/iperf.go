package main

import (
	"errors"
	"os/exec"
)

func (s *server) runIperfClient(testId int, addr, port string, onComplete func(string, bool) error) error {
	if s.process != nil {
		return errors.New("iPerf is already running")
	}
	s.log.Debug("Starting iPerf client with test ID ", testId, ", address ", addr, ", port ", port)

	s.testId = testId
	s.process = exec.Command("iperf", "-c", addr, "-p", port)
	go func() {
		out, err := s.process.CombinedOutput()
		if err != nil {
			s.log.Error("Test failed with error: " + err.Error())
			_ = onComplete(string(out), true)
			return
		}

		_ = onComplete(string(out), false)
	}()
	return nil
}

func (s *server) runIperfServer(testId int) error {
	if s.process != nil {
		return errors.New("iPerf is already running")
	}
	s.log.Debug("Starting iPerf server with test ID ", testId)

	s.testId = testId
	s.process = exec.Command("iperf", "-s", "-p", s.config.iperfPort)
	go func() {
		_, _ = s.process.CombinedOutput()
	}()
	return nil
}
func (s *server) closeIperf() error {
	if s.process == nil {
		return errors.New("iPerf is not running")
	}
	s.log.Debug("Closing iPerf with test ID ", s.testId)

	if err := s.process.Process.Kill(); err != nil {
		return err
	}
	s.process = nil

	return nil

}
