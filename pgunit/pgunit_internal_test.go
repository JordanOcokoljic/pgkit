package pgunit

import "testing"

func TestGenerateRandomName(t *testing.T) {
	for i := 0; i < 10; i++ {
		first := generateRandomName()
		second := generateRandomName()

		if first == second {
			t.Error("two generated names were identical")
			break
		}
	}
}
