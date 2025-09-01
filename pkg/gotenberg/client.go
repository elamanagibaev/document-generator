package gotenberg

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}
