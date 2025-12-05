package engine

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/Rishabhsingh78/orderMatchingEngine/pkg/utils"
)

func BenchmarkEngine_ProcessOrders(b *testing.B) {
	eng := NewEngine()
	symbol := "BTCUSD"

	// Pre-populate with some orders to simulate a realistic state
	for i := 0; i < 1000; i++ {
		eng.SubmitOrder(&Order{
			ID:        utils.GenerateUUID(),
			Symbol:    symbol,
			Side:      SideBuy,
			Type:      OrderTypeLimit,
			Price:     int64(50000 + rand.Intn(1000)),
			Quantity:  int64(1 + rand.Intn(10)),
			Timestamp: time.Now().UnixMilli(),
		})
		eng.SubmitOrder(&Order{
			ID:        utils.GenerateUUID(),
			Symbol:    symbol,
			Side:      SideSell,
			Type:      OrderTypeLimit,
			Price:     int64(51000 + rand.Intn(1000)),
			Quantity:  int64(1 + rand.Intn(10)),
			Timestamp: time.Now().UnixMilli(),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		side := SideBuy
		if i%2 == 0 {
			side = SideSell
		}

		price := int64(50500 + rand.Intn(1000)) // Overlapping price range to trigger matches

		eng.SubmitOrder(&Order{
			ID:        utils.GenerateUUID(),
			Symbol:    symbol,
			Side:      side,
			Type:      OrderTypeLimit,
			Price:     price,
			Quantity:  1,
			Timestamp: time.Now().UnixMilli(),
		})
	}
}

func TestLatency(t *testing.T) {
	eng := NewEngine()
	symbol := "ETHUSD"
	numOrders := 10000
	latencies := make([]int64, 0, numOrders)

	for i := 0; i < numOrders; i++ {
		side := SideBuy
		if i%2 == 0 {
			side = SideSell
		}
		price := int64(3000 + rand.Intn(100))

		order := &Order{
			ID:        utils.GenerateUUID(),
			Symbol:    symbol,
			Side:      side,
			Type:      OrderTypeLimit,
			Price:     price,
			Quantity:  1,
			Timestamp: time.Now().UnixMilli(),
		}

		start := time.Now()
		_, err := eng.SubmitOrder(order)
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Failed to submit order: %v", err)
		}

		latencies = append(latencies, duration.Microseconds())
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	p50 := latencies[int(float64(numOrders)*0.50)]
	p99 := latencies[int(float64(numOrders)*0.99)]
	p999 := latencies[int(float64(numOrders)*0.999)]

	fmt.Printf("\nLatency Results (microseconds):\n")
	fmt.Printf("p50: %d us\n", p50)
	fmt.Printf("p99: %d us\n", p99)
	fmt.Printf("p99.9: %d us\n", p999)
}
