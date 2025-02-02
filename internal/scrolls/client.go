package scrolls

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"

	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/library"
	"github.com/mreliasen/scrolls-cli/internal/settings"
)

type Client struct {
	baseUrl  *url.URL
	token    string
	version  string
	Settings *settings.Settings
	base     *client
	Library  *library.Library

	// Files   *FileClient
	Version *VersionClient
	Storage *StorageClient
}

type client struct {
	client *Client
}

func New() (*Client, error) {
	b, _ := url.Parse("https://get.scrolls.sh")

	sc := &Client{
		baseUrl: b,
		token:   "",
	}

	sc.base = &client{sc}

	// sc.Files = (*FileClient)(sc.base)
	sc.Storage = (*StorageClient)(sc.base)
	sc.Version = (*VersionClient)(sc.base)

	configSettings, err := settings.LoadSettings()
	if err != nil {
		return nil, fmt.Errorf("error reading settings file: %w", err)
	}

	sc.Settings = configSettings

	configLibrary, err := library.LoadLibrary()
	if err != nil {
		return nil, fmt.Errorf("error loading library: %w", err)
	}

	sc.Library = configLibrary

	return sc, nil
}

func (c *Client) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	url, err := url.Parse(c.baseUrl.String())
	if err != nil {
		return nil, err
	}

	url, err = url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Scrolls-CLI-Version", c.version)
	req.Header.Add("User-Agent", fmt.Sprintf("scrolls.sh-cli/%s (%s/%s)", c.version, runtime.GOOS, runtime.GOARCH))
	req.Header.Add("Content-Type", "Application/json")

	return req, nil
}

func (c *Client) apiCall(method, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(method, endpoint, body)

	if flags.Debug() {
		fmt.Println(req)
	}

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if flags.Debug() {
		fmt.Println(resp)
	}

	return resp, nil
}

func (c *Client) Get(path string, body io.Reader) (*http.Response, error) {
	return c.apiCall("GET", path, body)
}

func (c *Client) Post(path string, body io.Reader) (*http.Response, error) {
	return c.apiCall("POST", path, body)
}

func (c *Client) Patch(path string, body io.Reader) (*http.Response, error) {
	return c.apiCall("PATCH", path, body)
}

func (c *Client) Put(path string, body io.Reader) (*http.Response, error) {
	return c.apiCall("PUT", path, body)
}

func (c *Client) Delete(path string) (*http.Response, error) {
	return c.apiCall("DELETE", path, nil)
}
