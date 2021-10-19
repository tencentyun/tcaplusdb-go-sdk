package tcaplus

import "time"

type Client struct {
	*client
	defZone    int32
	defTimeout time.Duration
}

// 兼容老接口，保留
func NewClient() *Client {
	c := new(Client)
	c.client = newClient(false)
	c.defZone = -1
	c.defTimeout = 5 * time.Second
	return c
}

func NewTDRClient() *Client {
	c := new(Client)
	c.client = newClient(false)
	c.defZone = -1
	c.defTimeout = 5 * time.Second
	return c
}
