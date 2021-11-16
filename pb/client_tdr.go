package tcaplus

import "time"

type Client struct {
	*client
}

// 兼容老接口，保留
func NewClient() *Client {
	c := new(Client)
	c.client = newClient(false)
	c.defTimeout = 5 * time.Second
	return c
}

func NewTDRClient() *Client {
	c := new(Client)
	c.client = newClient(false)
	c.defTimeout = 5 * time.Second
	return c
}
