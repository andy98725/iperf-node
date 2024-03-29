package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (s *server) connect() error {
	body := struct {
		Id         int    `json:"id"`
		Key        string `json:"key"`
		ServerPort string `json:"serverPort"`
		IPerfPort  string `json:"iPerfPort"`
	}{
		Id:         s.config.id,
		Key:        s.config.key,
		ServerPort: s.config.hostPort,
		IPerfPort:  s.config.iperfPort,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.config.serverAddr+"/api/connect", &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			s.log.Fatal(err)
		}
		s.log.Error("Endpoint responded with code ", res.StatusCode)
		s.log.Error(string(bodyBytes))
		return errors.New(string(bodyBytes))
	}

	s.log.Info("Connected.")
	return nil
}

func (s *server) completeClientTest(results string, failed bool) error {
	if results == "" {
		results = "<no output>"
	}

	if failed {
		s.log.Error("Failed test.")
		s.log.Error(results)
	} else {
		s.log.Debug("Test completed.")
		s.log.Debug(results)
	}

	body := struct {
		Id      int    `json:"id"`
		Key     string `json:"key"`
		Results string `json:"results"`
		TestId  int    `json:"testId"`
		Failed  bool   `json:"failed"`
	}{
		Id:      s.config.id,
		Key:     s.config.key,
		Results: results,
		TestId:  s.testId,
		Failed:  failed,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.config.serverAddr+"/api/complete/client", &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			s.log.Fatal(err)
		}
		s.log.Error("Endpoint responded with code ", res.StatusCode)
		s.log.Error(string(bodyBytes))
	}

	return nil
}

func (s *server) completeServerTest() error {
	s.log.Debug("Completed test.")

	body := struct {
		Id     int    `json:"id"`
		Key    string `json:"key"`
		TestId int    `json:"testId"`
	}{
		Id:     s.config.id,
		Key:    s.config.key,
		TestId: s.testId,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.config.serverAddr+"/api/complete/server", &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			s.log.Error(err)
		}
		s.log.Error("Endpoint responded with code ", res.StatusCode)
		s.log.Error(string(bodyBytes))
	}

	return nil
}
