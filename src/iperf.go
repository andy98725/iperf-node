package main

import (
	"errors"
	"os/exec"
)

func (s *server) runIperfClient(testId int, addr, port string) error {
	if s.process != nil {
		return errors.New("iPerf is already running")
	}
	s.log.Debug("Starting iPerf client with test ID ", testId, ", address ", addr, ", port ", port)

	s.testId = testId
	s.process = exec.Command("iperf", "-c "+addr, "-p "+s.config.iperfPort)
	go func() {
		out, err := s.process.CombinedOutput()
		if err != nil {
			s.log.Error("Test failed with error: " + err.Error())
			s.failTest(string(out))
			return
		}

		s.completeClientTest(string(out))
	}()
	return nil
}

func (s *server) runIperfServer(testId int) error {
	if s.process != nil {
		return errors.New("iPerf is already running")
	}
	s.log.Debug("Starting iPerf server with test ID ", testId)

	s.testId = testId
	s.process = exec.Command("iperf", "-s", "-p "+s.config.iperfPort)
	go func() {
		_, _ = s.process.CombinedOutput()
	}()
	return nil
}
func (s *server) finishIperfServer() error {
	if s.process == nil {
		return errors.New("iPerf is not running")
	}
	s.log.Debug("Closing iPerf server with test ID ", s.testId)

	if err := s.process.Process.Kill(); err != nil {
		return err
	}
	s.process = nil

	s.completeServerTest()
	return nil
}
