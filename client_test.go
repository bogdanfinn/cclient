package cclient

import (
	"encoding/json"
	"github.com/Carcraftz/fhttp/http2"
	"io"
	"io/ioutil"
	"testing"

	tls "github.com/Carcraftz/utls"
)

type JA3Response struct {
	JA3Hash   string `json:"ja3_hash"`
	JA3       string `json:"ja3"`
	UserAgent string `json:"User-Agent"`
}

func readAndClose(r io.ReadCloser) ([]byte, error) {
	readBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return readBytes, r.Close()
}

const Chrome83Hash = "b32309a26951912be7dba376398abc3b"

var settings = map[http2.SettingID]uint32{
	http2.SettingHeaderTableSize:      65536,
	http2.SettingMaxConcurrentStreams: 1000,
	http2.SettingInitialWindowSize:    6291456,
	http2.SettingMaxHeaderListSize:    262144,
}

var settingsOrder = []http2.SettingID{
	http2.SettingHeaderTableSize,
	http2.SettingMaxConcurrentStreams,
	http2.SettingInitialWindowSize,
	http2.SettingMaxHeaderListSize,
}

var client, _ = NewClient(tls.HelloChrome_83, "", true, 6, settings, settingsOrder) // cannot throw an error because there is no proxy

func TestCClient_JA3(t *testing.T) {
	resp, err := client.Get("https://ja3er.com/json")
	if err != nil {
		t.Fatal(err)
	}

	respBody, err := readAndClose(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var ja3Response JA3Response
	if err := json.Unmarshal(respBody, &ja3Response); err != nil {
		t.Fatal(err)
	}

	if ja3Response.JA3Hash != Chrome83Hash {
		t.Error("unexpected JA3 hash; expected:", Chrome83Hash, "| got:", ja3Response.JA3Hash)
	}
}

func TestCClient_HTTP2(t *testing.T) {
	resp, err := client.Get("https://www.google.com")
	if err != nil {
		t.Fatal(err)
	}

	_, err = readAndClose(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.ProtoMajor != 2 || resp.ProtoMinor != 0 {
		t.Error("unexpected response proto; expected: HTTP/2.0 | got: ", resp.Proto)
	}
}
