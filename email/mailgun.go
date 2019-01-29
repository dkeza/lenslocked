package email

import (
	"fmt"
	"net/url"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to simplegallery.kezic.net"
	resetSubject   = "Instructions for resetting your password."
	resetBaseUrl   = "https://simplegallery.kezic.net/reset"
)

const welcomeText = `Hi there!

Welcome to simplegallery.kezic.net

Bye
`

const welcomeHTML = `Hi there!<br>
<br>
Welcome to simplegallery.kezic.net<br>
<br>
Bye<br>
`

const resetTextTmpl = `Hi there!

It appears that you have requested a password reset. If this was you, please folow link below to update your password:

%s

If you are asked for a token, please use the following value:

%s

If you didn't request a password reset you can safely ignore this email and your account will not be changed.

Best,
Simplegallery Support
`

const resetHTMLTmpl = `Hi there!<br/>
It appears that you have requested a password reset. If this was you, please folow link below to update your password:<br/>
<br/>
<a href="%s">%s</a><br/>
<br/>
If you are asked for a token, please use the following value:<br/>
<br/>
%s<br/>
<br/>
If you didn't request a password reset you can safely ignore this email and your account will not be changed.<br/>
<br/>
Best,<br/>
Simplegallery Support<br/>
`

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		c.mg = mg
	}
}

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		from: "support@simplegallery.kezic.net",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := mailgun.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseUrl + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	message := mailgun.NewMessage(c.from, resetSubject, resetText, toEmail)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, resetUrl, token)
	message.SetHtml(resetHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%v <%v>", name, email)
}
