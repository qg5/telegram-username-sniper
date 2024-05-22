package telegram

import (
	"context"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
)

type Client struct {
	client *telegram.Client
	api    *tg.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new Telegram client, handles authentication, and runs it in a background goroutine.
func New(appID int, appHash, phoneNumber string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	client := telegram.NewClient(appID, appHash, telegram.Options{})
	tgClient := &Client{
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}

	passedAuthFlow := make(chan struct{})
	authFlow := auth.NewFlow(SimpleAuthFlow{PhoneNumber: phoneNumber}, auth.SendCodeOptions{})
	
	go func() {
		err := client.Run(ctx, func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, authFlow); err != nil {
				return err
			}

			tgClient.api = client.API()
			close(passedAuthFlow)

			<-ctx.Done()
			return ctx.Err()
		})

		if err != nil {
			duration, ok := tgerr.AsFloodWait(err)
			if ok {
				log.Fatalf("Flood wait hit, cant signin for: %v\n", duration)
				return
			}

			log.Fatalf("Couldnt run client: %v\n", err)
		}
	}()

	<-passedAuthFlow
	return tgClient
}

func (c *Client) CreateChannel(accessHash string) error {
	u, err := c.api.ChannelsCreateChannel(context.Background(), &tg.ChannelsCreateChannelRequest{Title: accessHash, Broadcast: true})
	if err != nil {
		return err
	}

	inputChannel := &tg.InputChannel{}

	switch u := u.(type) {
	case *tg.Updates:
		channel := u.Chats[0].(*tg.Channel)
		inputChannel = &tg.InputChannel{ChannelID: channel.GetID(), AccessHash: channel.AccessHash}
	}

	if _, err := c.api.ChannelsUpdateUsername(c.ctx, &tg.ChannelsUpdateUsernameRequest{
		Channel:  inputChannel,
		Username: accessHash,
	}); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateUsername(newUsername string) error {
	_, err := c.api.AccountUpdateUsername(c.ctx, newUsername)
	return err
}
