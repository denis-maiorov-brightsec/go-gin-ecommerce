package http

import "testing"

func TestParseOrderIDAcceptsPositiveInteger(t *testing.T) {
	id, err := ParseOrderID("42")
	if err != nil {
		t.Fatalf("expected valid order id, got %v", err)
	}

	if id != 42 {
		t.Fatalf("expected order id 42, got %d", id)
	}
}

func TestParseOrderIDRejectsInvalidValues(t *testing.T) {
	for _, rawID := range []string{"", "0", "-1", "abc"} {
		if _, err := ParseOrderID(rawID); err == nil {
			t.Fatalf("expected %q to be rejected", rawID)
		}
	}
}
