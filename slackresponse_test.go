package slackresponse

import (
	"fmt"
	"log"
	"os"
	"testing"

	jsonextract "github.com/blainemoser/JsonExtract"
	logging "github.com/blainemoser/Logging"
	utils "github.com/blainemoser/goutils"
)

var (
	l   *logging.Log
	dir string
)

func TestMain(m *testing.M) {
	var err error
	dir, err = utils.BaseDir([]string{"api"}, "FinWatch")
	if err != nil {
		log.Fatal(err)
	}
	l, err = logging.NewLog(fmt.Sprintf("%s/log.log", dir), "TEST")
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestPost(t *testing.T) {
	url, err := getTestURL()
	if err != nil {
		t.Fatal(err)
	}
	if len(url) < 1 {
		t.Fatalf("could not complete test, missing url")
	}
	err = SlackPost("Test Post", "this is some testing content\n this should post to slack", "INFO", url, l)
	if err != nil {
		t.Error(err)
	}
}

func getTestURL() (string, error) {
	content, err := utils.GetFileContent(fmt.Sprintf("%s/%s", dir, "test_env.json"))
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("test environment not found: %s", err.Error())
			return "", nil
		}
		return "", err
	}
	return getURL(content)
}

func getURL(content []byte) (string, error) {
	js := jsonextract.JSONExtract{RawJSON: string(content)}
	if URLinterface, err := js.Extract("slackURL"); err == nil {
		if url, ok := URLinterface.(string); ok {
			return url, nil
		} else {
			return "", fmt.Errorf("could not parse slack url")
		}
	} else {
		return "", err
	}
}
