package chatbase

import "testing"

func TestTimeStamp(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		result := TimeStamp()
		if result <= 0 {
			t.Errorf("Expected non-zero timestamp, got %v", result)
		}
	})
}
