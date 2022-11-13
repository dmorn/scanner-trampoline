package config

import (
	"testing"
)

func TestUmarshal(t *testing.T) {
	doc := `
	open_cmd = ["open"]
	`

	var cfg Config
	if err := Unmarshal([]byte(doc), &cfg); err != nil {
		t.Fatal(err)
	}

	expectedOpen := []string{"open"}
	if len(expectedOpen) != len(cfg.OpenCmd) {
		t.Fatalf("expected len %d, have %d", len(expectedOpen), len(cfg.OpenCmd))
	}
	for i, have := range cfg.OpenCmd {
		want := expectedOpen[i]
		if have != want {
			t.Fatalf("wanted %v, have %v", want, have)
		}
	}
}
