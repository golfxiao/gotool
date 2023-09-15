package waitgroup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimitedWaitGroup_Add(t *testing.T) {
	// Test case 1: Adding one ticket to the LimitedWaitGroup
	lwg := NewLimitedWaitGroup(4)
	lwg.Add(1)
	if len(lwg.tickets) != 1 {
		t.Errorf("Expected 1 ticket, got %d", len(lwg.tickets))
	}

	// Test case 2: Adding multiple tickets to the LimitedWaitGroup
	lwg.Add(3)
	if len(lwg.tickets) != 4 {
		t.Errorf("Expected 4 tickets, got %d", len(lwg.tickets))
	}
}

func TestLimitedWaitGroup_Done(t *testing.T) {
	wg := NewLimitedWaitGroup(1)

	// Test case 1: Ensure that Done() releases a ticket and calls wg.Done()
	wg.Add(1) // Add a task to the wait group
	wg.Done() // Release the ticket

	// Test case 2: Ensure that Done() works correctly when called multiple times
	wg.Add(1) // Add a task to the wait group
	wg.Done() // Release the ticket
	wg.Done() // Release the ticket again
}

func TestLimitedWaitGroup_Wait(t *testing.T) {
	tasks, limit := 100, 10
	startTime := time.Now()

	wg := NewLimitedWaitGroup(limit)

	for i := 0; i < tasks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Second * 1)
		}()
	}

	wg.Wait()

	assert.Equal(t, int(tasks/limit), int(time.Since(startTime)/time.Second))
}
