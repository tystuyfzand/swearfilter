package swearfilter

import (
	"testing"
)

func TestNew(t *testing.T) {
	filter := New([]string{"foo", "bar"}).(*swearFilter)
	if filter.disableNormalize {
		t.Errorf("Filter option disableNormalize was incorrect, got: %t, want: %t", filter.disableNormalize, false)
	}
	if filter.disableSpacedTab {
		t.Errorf("Filter option disableSpacedTab was incorrect, got: %t, want: %t", filter.disableSpacedTab, false)
	}
	if filter.disableMultiWhitespaceStripping {
		t.Errorf("Filter option disableMultiWhitespaceStripping was incorrect, got: %t, want: %t", filter.disableMultiWhitespaceStripping, false)
	}
	if filter.disableZeroWidthStripping {
		t.Errorf("Filter option disableZeroWidthStripping was incorrect, got: %t, want: %t", filter.disableZeroWidthStripping, false)
	}
	if !filter.enableSpacedBypass {
		t.Errorf("Filter option enableSpacedBypass was incorrect, got: %t, want: %t", filter.enableSpacedBypass, true)
	}
	if len(filter.BadWords) != 2 {
		t.Errorf("Filter option BadWords was incorrect, got length: %d, want length: %d", len(filter.BadWords), 2)
	}
}

func TestCheck(t *testing.T) {
	filter := New([]string{"foo"}, EnableSpacedBypass())
	messages := []string{"fooing", "foo", "asdf", "what the f o o dude"}

	for i := 0; i < len(messages); i++ {
		trippers, err := filter.Check(messages[i])
		if err != nil {
			t.Errorf("Check failed due to external dependency: %v", err)
		}
		switch i {
		case 0, 1, 3:
			if len(trippers) != 1 {
				t.Fatalf("Check %d did not act as expected, got trippers length: %d, want trippers length: %d", i, len(trippers), 1)
			}
			if trippers[0].Word != "foo" {
				t.Errorf("Check did not act as expected, got first tripper: %s, want first tripper: %s", trippers[0].Word, "foo")
			}
		case 2:
			if len(trippers) != 0 {
				t.Errorf("Check did not act as expected, got trippers length: %d, want trippers length: %d", len(trippers), 0)
			}
		default:
			t.Errorf("Check test invalid, got test messages length: %d, want test messages length: %d", len(messages), 4)
		}
	}
}
