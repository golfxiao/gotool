package weighted

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDraw(t *testing.T) {
	weightedRandom := NewWeightedRandom(50)

	// Testing when winCount is empty
	winCount := []int{}
	num := 5
	expected := []int{}
	result := weightedRandom.Draw(winCount, num)
	assert.Equal(t, expected, result)

	// Testing when num is negative
	winCount = []int{1, 2, 3}
	num = -3
	expected = []int{1}
	result = weightedRandom.Draw(winCount, num)
	assert.Equal(t, len(expected), len(result))

	// Testing when num is greater than the length of winCount
	winCount = []int{1, 2, 3}
	num = 5
	expected = []int{0, 1, 2}
	result = weightedRandom.Draw(winCount, num)
	sort.Ints(result)
	assert.Equal(t, expected, result)
}

func TestDrawOne(t *testing.T) {
	n := 100
	winCount := make([]int, 100)
	w := NewWeightedRandom(50)
	for i := 0; i < n; i++ {
		u := w.DrawOne(winCount)
		winCount[u]++
	}

	t.Logf("The distribution data of winner users: %v", winCount)
}
