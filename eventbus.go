package dgrs

import "context"

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
				return s.options.UnmarshalFunc([]byte(msg.String()), v)
			})
		}
	}()
	return
}

func (s *State) joinChanKeys(channel string) string {
	return s.joinKeys("chan", channel)
}
