package facebook

import (
	"encoding/json"
	"net/http"
)

// Client provides Facebook Graph API accessing
type Client struct {
	accessToken string
	client      *http.Client
}

// NewClient returns a new Client instance
func NewClient(appAccessToken string, transport http.RoundTripper) Client {
	return Client{
		accessToken: appAccessToken,
		client:      &http.Client{Transport: transport},
	}
}

// GetURLNodeFields gets fields on a URLNode.
func (c Client) GetURLNodeFields(id string) (URLNode, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://graph.facebook.com/v6.0/", nil)
	q := req.URL.Query()
	q.Set("id", id)
	q.Set("access_token", c.accessToken)
	q.Set("fields", "engagement")
	req.URL.RawQuery = q.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return URLNode{}, err
	}
	defer res.Body.Close()

	var un URLNode
	err = json.NewDecoder(res.Body).Decode(&un)

	return un, err
}
