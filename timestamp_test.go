package chatbase

import "testing"

func TestTimeStamp(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		result := TimeStamp()
		if result < 10000000000 {
			t.Errorf("Expected non-zero timestamp in milliseconds, got %v", result)
		}
	})
}
