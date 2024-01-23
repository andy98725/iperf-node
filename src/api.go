package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (s *server) connect() error {
	body := struct {
		Id         int
		Key        string
		ServerPort string
		iPerfPort  string
	}{
		Id:         s.config.id,
		Key:        s.config.key,
		ServerPort: s.config.hostPort,
		iPerfPort:  s.config.iperfPort,
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

	resp := &struct{ Status string }{}
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return err
	}

	if res.StatusCode != 200 {
		s.log.Debug("response recieved:")
		s.log.Debug(resp.Status)
	}

	return nil
}

func (s *server) completeTest(results *string) error {
	var buf bytes.Buffer

	if results != nil {
		body := struct {
			Id      int
			Key     string
			Results string
			TestId  int
		}{
			Id:      s.config.id,
			Key:     s.config.key,
			Results: *results,
			TestId:  s.testId,
		}
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
	} else {
		body := struct {
			Id     int
			Key    string
			TestId int
		}{
			Id:     s.config.id,
			Key:    s.config.key,
			TestId: s.testId,
		}
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest("POST", s.config.serverAddr+"/api/complete", &buf)
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

	resp := &struct{ Status string }{}
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return err
	}

	if res.StatusCode != 200 {
		s.log.Debug("response recieved:")
		s.log.Debug(resp.Status)
	}

	return nil
}
