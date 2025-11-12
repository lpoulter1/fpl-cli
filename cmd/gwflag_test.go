package cmd

import (
	"testing"
)

func TestGWFlagParsing(t *testing.T) {
	var flag gwFlag
	if err := flag.Set("1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(flag.Ranges()) != 1 || flag.Ranges()[0] != (gwRange{Start: 1, End: 1}) {
		t.Fatalf("unexpected ranges: %+v", flag.Ranges())
	}

	flag = gwFlag{}
	if err := flag.Set("1-3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := flag.Ranges()[0]; got.Start != 1 || got.End != 3 {
		t.Fatalf("expected 1-3 got %+v", got)
	}

	flag = gwFlag{}
	if err := flag.Set("1|3-4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := flag.Set("2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(flag.Ranges()) != 1 || flag.Ranges()[0] != (gwRange{Start: 1, End: 4}) {
		t.Fatalf("expected merged range 1-4, got %+v", flag.Ranges())
	}
}

func TestGWFlagInvalid(t *testing.T) {
	var flag gwFlag
	if err := flag.Set("0"); err == nil {
		t.Fatal("expected error for gw 0")
	}

	if err := flag.Set("3-1"); err == nil {
		t.Fatal("expected error for inverted range")
	}
}
