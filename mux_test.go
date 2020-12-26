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
	{[]string{"bob.notification", "alice.trades", "john.products"}, "bob.notification,alice.trades,john.products", true},
	{[]string{"bob.notification", "john.products"}, "bob.notification,alice.trades,john.products", false},
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
	{[]string{"bob.notification", "alice.trades", "john.products"}, "bob.notification,alice.trades,john.products", 0},
	{[]string{"bob.notification", "john.products"}, "bob.notification,alice.trades,john.products", 1},
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
