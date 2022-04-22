package cclient

import (
	"github.com/Carcraftz/fhttp/http2"
	"time"

	http "github.com/Carcraftz/fhttp"
	"github.com/Carcraftz/fhttp/cookiejar"

	"golang.org/x/net/proxy"

	utls "github.com/Carcraftz/utls"
)

func NewClient(clientHello utls.ClientHelloID, proxyUrl string, allowRedirect bool, timeout time.Duration, settings map[http2.SettingID]uint32, settingsOrder []http2.SettingID, insecureSkipVerify bool) (http.Client, error) {
	if len(proxyUrl) > 0 {
		dialer, err := newConnectDialer(proxyUrl)
		if err != nil {
			if allowRedirect {
				cJar, _ := cookiejar.New(nil)
				return http.Client{
					Jar:     cJar,
					Timeout: time.Second * timeout,
				}, err
			}
			cJar, _ := cookiejar.New(nil)
			return http.Client{
				Jar:     cJar,
				Timeout: time.Second * timeout,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}, err
		}
		if allowRedirect {
			cJar, _ := cookiejar.New(nil)
			return http.Client{
				Jar:       cJar,
				Transport: newRoundTripper(clientHello, settings, settingsOrder, insecureSkipVerify, dialer),
				Timeout:   time.Second * timeout,
			}, nil
		}
		cJar, _ := cookiejar.New(nil)
		return http.Client{
			Jar:       cJar,
			Transport: newRoundTripper(clientHello, settings, settingsOrder, insecureSkipVerify, dialer),
			Timeout:   time.Second * timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}, nil
	} else {
		if allowRedirect {
			cJar, _ := cookiejar.New(nil)
			return http.Client{
				Jar:       cJar,
				Transport: newRoundTripper(clientHello, settings, settingsOrder, insecureSkipVerify, proxy.Direct),
				Timeout:   time.Second * timeout,
			}, nil
		}
		cJar, _ := cookiejar.New(nil)
		return http.Client{
			Jar:       cJar,
			Transport: newRoundTripper(clientHello, settings, settingsOrder, insecureSkipVerify, proxy.Direct),
			Timeout:   time.Second * timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}, nil

	}
}
