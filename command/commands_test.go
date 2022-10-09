package command

import (
	"fmt"
	"testing"
)

const CMDCLOSE = "close"
const CMDADD = "add"

func TestManager_Parse_NoArg(t *testing.T) {
	var closeRun = false

	manager := NewManager()
	manager.Add(
		NewItem(CMDCLOSE),
		func(strings []string) error {
			closeRun = true
			return nil
		},
	)

	f := manager.Parse([]byte(CMDCLOSE))
	if err := f(); err != nil {
		t.Fatalf("No error expected: %v", err)
	}

	if !closeRun {
		t.Errorf("closeRun not run")
	}
}

func TestManager_Parse_WithOneArg(t *testing.T) {
	var addRun = false

	manager := NewManager()
	manager.Add(
		NewItem(CMDADD),
		func(strings []string) error {
			if len(strings) != 1 {
				t.Fatalf("should receive one arg, got %d: %v", len(strings), strings)
			}
			if strings[0] != "10" {
				t.Fatalf("Should receive 10")
			}
			addRun = true
			return nil
		},
	)

	f := manager.Parse([]byte(fmt.Sprintf("%s 10", CMDADD)))
	if err := f(); err != nil {
		t.Fatalf("No error expected: %v", err)
	}

	if !addRun {
		t.Errorf("addRun not run")
	}
}

func TestManager_Parse_With4Arg(t *testing.T) {
	var addRun = false

	manager := NewManager()
	manager.Add(
		NewItem(CMDADD),
		func(strings []string) error {
			if len(strings) != 4 {
				t.Fatalf("should receive one arg, got %d: %v", len(strings), strings)
			}
			if strings[0] != "10" {
				t.Fatalf("Should receive 10")
			}
			if strings[1] != "test" {
				t.Fatalf("Should receive test")
			}
			if strings[2] != "84848.333" {
				t.Fatalf("Should receive 84848.333")
			}
			if strings[3] != "39" {
				t.Fatalf("Should receive 39")
			}
			addRun = true
			return nil
		},
	)

	f := manager.Parse([]byte(fmt.Sprintf("%s 10 test 84848.333 39", CMDADD)))
	if err := f(); err != nil {
		t.Fatalf("No error expected: %v", err)
	}

	if !addRun {
		t.Errorf("addRun not run")
	}
}
