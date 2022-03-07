package dgrs

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

const dmBcChannel = "internal-dm-broadcasts"

type DirectMessageEvent struct {
	Message  *discordgo.Message `json:"message"`
	Channel  *discordgo.Channel `json:"channel"`
	IsUpdate bool               `json:"is_update"`
}

// Publish pushes an event payload to the given
// channel to all subscribers connected to the
// same Redis instance.
func (s *State) Publish(channel string, payload interface{}) (err error) {
	data, err := s.options.MarshalFunc(payload)
	if err != nil {
		return
	}
	ctx, cf := s.getContext()
	defer cf()
	err = s.client.Publish(ctx, s.joinChanKeys(channel), data).Err()
	return
}

// Subscribe starts an event listener on the given
// channel and passes all received events to the
// passed handler function. The handler is getting
// passed a scan function which takes an object
// instance to scan the payload into.
//
// Subscribe returns a close function to close
// the listener on the channel.
func (s *State) Subscribe(
	channel string,
	handler func(scan func(v interface{}) error),
) (close func() error) {
	pubsub := s.client.Subscribe(context.Background(), s.joinChanKeys(channel))
	ch := pubsub.Channel()
	close = func() error {
		return pubsub.Close()
	}
	go func() {
		for msg := range ch {
			go handler(func(v interface{}) error {
				return s.options.UnmarshalFunc([]byte(msg.Payload), v)
			})
		}
	}()
	return
}

// SubscribeDMs subscribes to the states DM event bus
// and executes handler on each received DM event with the
// event details passed as payload.
func (s *State) SubscribeDMs(handler func(e *DirectMessageEvent)) (close func() error) {
	return s.Subscribe(dmBcChannel, func(scan func(v interface{}) error) {
		var dm DirectMessageEvent
		if err := scan(&dm); err == nil {
			handler(&dm)
		}
	})
}

func (s *State) publishDM(msg *discordgo.Message, isUpdate bool) {
	if !s.options.BroadcastDMs {
		return
	}
	ch, err := s.Channel(msg.ChannelID)
	if err != nil || ch == nil {
		return
	}
	if ch.Type != discordgo.ChannelTypeDM && ch.Type != discordgo.ChannelTypeGroupDM {
		return
	}
	s.Publish(dmBcChannel, DirectMessageEvent{msg, ch, isUpdate})
}

func (s *State) joinChanKeys(channel string) string {
	return s.joinKeys("chan", channel)
}
