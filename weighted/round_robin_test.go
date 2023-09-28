package weighted

import "testing"

func TestRoundRobin(t *testing.T) {
	servers := []string{"server1", "server2", "server3"}
	selector := Selector{servers: servers, seq: 0}

	// Testing the first round-robin iteration
	expected := "server1"
	result := selector.RoundRobin()
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Testing the second round-robin iteration
	expected = "server2"
	result = selector.RoundRobin()
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Testing the third round-robin iteration
	expected = "server3"
	result = selector.RoundRobin()
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Testing the fourth round-robin iteration (should loop back to the first server)
	expected = "server1"
	result = selector.RoundRobin()
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
