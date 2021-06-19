package dgrc

import (
	"context"
	"strings"
)

func joinKeys(keys ...string) string {
	n := len(keyPrefix) + len(keys)
	for i := 0; i < len(keys); i++ {
		n += len(keys[i])
	}

	b := strings.Builder{}
	b.Grow(n)
	b.WriteString(keyPrefix)
	for _, s := range keys {
		b.WriteRune(keySeperator)
		b.WriteString(s)
	}

	return b.String()
}

func (s *State) getContext() (ctx context.Context) {
	if s.options.CommandTimeout > 0 {
		ctx, _ = context.WithTimeout(context.Background(), s.options.CommandTimeout)
	} else {
		ctx = context.Background()
	}
	return
}
