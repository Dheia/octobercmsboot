package octobercmsboot

import (
	"testing"
)

func TestOctober(t *testing.T) {

	t.Run("environment length", func(t *testing.T) {
		octobercms, _ := NewOctober("./template/OctoberConf.yaml", "dev")
		got := len(octobercms.Env)
		want := 2
		assertEquals(t, got, want)

	})
	t.Run("current environment", func(t *testing.T) {
		envTests := []struct {
			env string
		}{
			{env: "prod"},
			{env: "staging"},
			{env: "dev"},
		}
		for _, tt := range envTests {
			t.Run(tt.env, func(t *testing.T) {
				octobercms, _ := NewOctober("./template/OctoberConf.yaml", tt.env)
				got := octobercms.currentEnv
				want := tt.env
				assertEquals(t, got, want)
			})
		}
	})
}

func assertEquals(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
