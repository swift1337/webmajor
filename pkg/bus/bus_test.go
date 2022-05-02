package bus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBus_Publish(t *testing.T) {
	t.Run("PublishSimple", func(t *testing.T) {
		b := New(5)

		result := make(chan any, 2)

		s1 := b.Subscribe()
		s2 := b.Subscribe()

		b.Publish("hello")

		result <- s1.Channel()
		result <- s2.Channel()

		assert.Equal(t, len(result), 2)
	})

	t.Run("PublishSkipBusyChannel", func(t *testing.T) {
		b := New(1)

		s1 := b.Subscribe()

		b.Publish("hello")
		b.Publish("hello")

		assert.Equal(t, len(s1.Channel()), 1)
	})
}

func TestBus_Unsubscribe(t *testing.T) {
	t.Run("Unsubscribe", func(t *testing.T) {
		b := New(5)

		result := make(chan any, 2)

		s1 := b.Subscribe()

		result <- s1.Channel()

		b.Unsubscribe(s1)

		assert.Equal(t, b.Subscribers(), 0)
	})

	t.Run("UnsubscribeTwiceHandled", func(t *testing.T) {
		b := New(5)

		result := make(chan any, 2)

		s1 := b.Subscribe()

		result <- s1.Channel()

		b.Unsubscribe(s1)
		b.Unsubscribe(s1)

		assert.Equal(t, b.Subscribers(), 0)
	})
}
