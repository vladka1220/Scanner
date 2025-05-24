package notifier

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestEscapeMarkdown(t *testing.T) {
	input := "special _chars* to [escape]"
	expected := "special \\_chars\\* to \\[escape\\]"
	if out := escapeMarkdown(input); out != expected {
		t.Fatalf("unexpected escape result: %s", out)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSendTelegramMessage(t *testing.T) {
	os.Setenv("TELEGRAM_BOT_TOKEN", "tok")
	os.Setenv("TELEGRAM_CHAT_ID", "123")

	var captured url.Values
	rt := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "api.telegram.org" {
			t.Fatalf("unexpected host: %s", req.URL.Host)
		}
		body, _ := ioutil.ReadAll(req.Body)
		captured, _ = url.ParseQuery(string(body))
		return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(nil))}, nil
	})

	oldClient := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: rt}
	defer func() { http.DefaultClient = oldClient }()

	if err := SendTelegramMessage("hi"); err != nil {
		t.Fatalf("SendTelegramMessage error: %v", err)
	}
	if captured.Get("chat_id") != "123" || captured.Get("text") != "hi" {
		t.Fatalf("unexpected form: %v", captured)
	}
}
