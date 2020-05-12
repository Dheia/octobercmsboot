package octobercmsboot

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	env := Env{}
	got := len(env.generateKey())
	want := 32
	if got != 32 {
		t.Errorf("got %q want %q", got, want)
	}
}
