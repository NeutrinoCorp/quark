package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicBool(t *testing.T) {
	t.Run("Verify atomicity of atomic bool", func(t *testing.T) {
		var b atomicBool = 0
		assert.Equal(t, false, b.isSet())
		go func() {
			b.setTrue()
		}()
		go func() {
			b.setTrue()
		}()
		assert.Equal(t, false, b.isSet())
	})
}

func BenchmarkAtomicBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Run("Stress out atomicity of atomic bool", func(b *testing.B) {
			b.ReportAllocs()
			var bl atomicBool = 0
			bl.setTrue()
			bl.setFalse()
			bl.isSet()
		})
	}
}
