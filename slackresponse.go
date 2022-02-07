package slackresponse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	logging "github.com/blainemoser/Logging"
)

type Slack struct {
	url string
	log *logging.Log
}

const blankPayload = `{
"text": "%s",
"attachments": [
		{
			"mrkdwn_in": ["text"],
			"color": "%s",
			"fields": %s,
			"ts": %d
		}
	]
}`

func SlackPost(heading, content, level, url string, log *logging.Log) error {
	s := &Slack{
		url: url,
		log: log,
	}
	data, err := s.getData(heading, content, level)
	if err != nil {
		return err
	}
	result, err := s.slackPost(data)
	if err != nil {
		return err
	}
	if result < 200 || result > 299 {
		err = fmt.Errorf("responded with a %d error: %s", result, http.StatusText(result))
	}
	return err
}

func (s *Slack) getData(heading, content, level string) (data []byte, err error) {
	var fields string
	fields, err = s.fields(content)
	if err != nil {
		return data, err
	}
	data = []byte(fmt.Sprintf(blankPayload, heading, s.getColor(level), fields, time.Now().Unix()))
	return data, err
}

func (s *Slack) getColor(level string) string {
	switch strings.ToUpper(level) {
	case "INFO":
		return "#ccf1fb"
	case "ERROR":
		return "#e2475b"
	case "CRITICAL":
		return "#e2475b"
	case "DEBUG":
		return "#ecb230"
	case "WARNING":
		return "#ecb230"
	default:
		return "#ccf1fb"
	}
}

func (s *Slack) slackPost(data []byte) (int, error) {
	s.log.Write(fmt.Sprintf("posting to slack: %s", string(data)), "INFO")
	resp, err := http.Post(s.url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func (s *Slack) fields(content string) (string, error) {
	contentSplit := strings.Split(content, "\n")
	result := make([]interface{}, 0)
	for _, line := range contentSplit {
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		result = append(result, map[string]string{
			"value": line,
		})
	}
	ret, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}
