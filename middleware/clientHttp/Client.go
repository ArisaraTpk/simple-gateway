package clientHttp

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog"
	"io"
	"time"
)

type Client struct {
	*req.Client
	Config *ClientConfig
	logger zerolog.Logger
}

func NewClient(cfg *ClientConfig) *Client {
	c := Client{
		Config: cfg,
	}
	c.Client = c.NewReqClient(cfg.BaseUrl)

	return &c
}

func (c *Client) NewReqClient(baseUrl string) *req.Client {
	return req.C().
		SetBaseURL(baseUrl).
		SetCommonErrorResult(&ErrorMessage{}).
		OnBeforeRequest(c.OnBeforeRequest).
		OnAfterResponse(c.OnAfterRequest)
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func (msg *ErrorMessage) Error() string {
	return fmt.Sprintf("API Error: %s", msg.Message)
}

func (c *Client) OnBeforeRequest(client *req.Client, req *req.Request) error {
	c.logger.Info().
		Str("body", string(req.Body)).
		Str("header", req.HeaderToString()).
		Msg(fmt.Sprintf("Client Request %s: %s", req.Method, req.RawURL))

	return nil

}

func (c *Client) OnAfterRequest(client *req.Client, resp *req.Response) error {
	if resp.Err != nil {
		return nil
	}
	if errMsg, ok := resp.ErrorResult().(*ErrorMessage); ok {
		resp.Err = errMsg
		return nil
	}
	if !resp.IsSuccessState() {
		resp.Err = fmt.Errorf("bad status: %s\nraw content:\n%s", resp.Status, resp.Dump())
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Err(err)
	}
	c.logger.Info().
		Str("body", string(b)).
		Str("timeResponse", time.Now().Sub(resp.Request.StartTime).String()).
		Str("header", resp.HeaderToString()).
		Msg(fmt.Sprintf("Client Response %s: %s", resp.Request.Method, resp.Request.RawURL))
	return nil
}

func (c *Client) SetLogger(logger zerolog.Logger) {
	c.logger = logger
}
