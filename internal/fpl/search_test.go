package fpl

import "testing"

func TestFindPlayerByName(t *testing.T) {
	players := []Element{
		{ID: 1, WebName: "Haaland", FirstName: "Erling", SecondName: "Haaland"},
		{ID: 2, WebName: "Haalan", FirstName: "Erik", SecondName: "Haalan"},
	}

	player, suggestions, err := FindPlayerByName("erling", players)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player == nil || player.ID != 1 {
		t.Fatalf("expected player 1, got %+v", player)
	}
	if len(suggestions) == 0 {
		t.Fatal("expected suggestions for fuzzy search")
	}
}

func TestFindPlayerByNameEmpty(t *testing.T) {
	if _, _, err := FindPlayerByName(" ", nil); err == nil {
		t.Fatal("expected error on empty name")
	}
}
