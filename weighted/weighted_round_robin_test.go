package weighted

import "testing"

func TestNewWeightedSelector(t *testing.T) {
	// Test case 1: Empty server map
	servers := make(map[string]int)
	selector := NewWeightedSelector(servers)
	if len(selector.servers) != 0 {
		t.Errorf("Expected 0 servers, got %d", len(selector.servers))
	}

	// Test case 2: Non-empty server map
	servers = map[string]int{
		"server1": 10,
		"server2": 5,
		"server3": 7,
	}
	selector = NewWeightedSelector(servers)
	if len(selector.servers) != len(servers) {
		t.Errorf("Expected %d servers, got %d", len(servers), len(selector.servers))
	}

	// Test case 3: Check totalWeight
	totalWeight := 0
	for _, weight := range servers {
		totalWeight += weight
	}
	if totalWeight != selector.totalWeight {
		t.Errorf("Expected totalWeight to be %d, got %d", totalWeight, selector.totalWeight)
	}
}
