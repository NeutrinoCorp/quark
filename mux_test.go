package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMux(t *testing.T) {
	t.Run("New mux", func(t *testing.T) {
		mux := NewMux()
		assert.NotNil(t, mux)
	})
}

var defaultMuxAddTestingSuite = []struct {
	c        *Consumer
	expected int
}{
	{c: &Consumer{topics: []string{"alex.trades", "bob.trades"}}, expected: 2},
	{c: &Consumer{topics: []string{"alex.trades"}}, expected: 1},
	{c: nil, expected: 0},
}

func TestDefaultMux_Add(t *testing.T) {
	for _, tt := range defaultMuxAddTestingSuite {
		t.Run("Default mux add", func(t *testing.T) {
			mux := NewMux()
			mux.Add(tt.c)
			assert.Equal(t, tt.expected, len(mux.List()))
		})
	}
}

var defaultMuxContainsTestingSuite = []struct {
	t        []string
	f        string
	expected bool
}{
	{[]string{"bob.notification", "alice.trades", "john.products"}, "alice.trades", true},
	{[]string{"bob.notification", "john.products"}, "alice.trades", false},
	{[]string{}, "alice.trades", false},
}

func TestDefaultMux_Contains(t *testing.T) {
	for _, tt := range defaultMuxContainsTestingSuite {
		t.Run("Default mux contains", func(t *testing.T) {
			mux := NewMux()
			mux.Topics(tt.t...)
			assert.Equal(t, tt.expected, mux.Contains(tt.f))
		})
	}
}

var defaultMuxDelTestingSuite = []struct {
	t        []string
	d        string
	expected int
}{
	{[]string{"bob.notification", "alice.trades", "john.products"}, "alice.trades", 2},
	{[]string{"bob.notification", "john.products"}, "alice.trades", 2},
	{[]string{}, "alice.trades", 0},
}

func TestDefaultMux_Del(t *testing.T) {
	for _, tt := range defaultMuxDelTestingSuite {
		t.Run("Default mux delete", func(t *testing.T) {
			mux := NewMux()
			mux.Topics(tt.t...)
			mux.Del(tt.d)
			assert.Equal(t, tt.expected, len(mux.List()))
		})
	}
}

func BenchmarkMux(b *testing.B) {
	mux := NewMux()
	mux.Topics("foo", "bar", "baz")
	b.Run("New Event Mux", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = NewMux()
		}
	})
	b.Run("Event Mux data insertion", func(b *testing.B) {
		b.ReportAllocs()
		mux := NewMux()
		for i := 0; i < b.N; i++ {
			mux.Topics("foo", "bar", "baz")
		}
	})
	b.Run("Event Mux single data insertion", func(b *testing.B) {
		b.ReportAllocs()
		mux := NewMux()
		for i := 0; i < b.N; i++ {
			mux.Topic("foo")
		}
	})
	b.Run("Event Mux data removal", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			mux.Del("baz")
		}
	})
	b.Run("Event Mux contains item", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			mux.Contains("baz")
		}
	})
	b.Run("Event Mux data manipulation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			mux.Get("baz")
			mux.Contains("bar")
			mux.List()
			mux.Del("baz")
		}
	})
}
