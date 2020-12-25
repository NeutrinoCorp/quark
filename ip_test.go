package quark

import "testing"

func BenchmarkIp(b *testing.B) {
	b.Run("Get local IP address", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Not used anymore, left to keep code coverage
			// _ = getLocalIP()
		}
	})
}
