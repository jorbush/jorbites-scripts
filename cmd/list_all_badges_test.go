package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAllBadgesCommand(t *testing.T) {
	t.Run("Success path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]string{"level_100", "week_streak"})
		}))
		defer server.Close()
		SetAppURL(server.URL)
		defer SetAppURL("http://localhost:3000")
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"list-all-badges"})
		err := RootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Available Badges in Jorbites (total 2):")
		assert.Contains(t, buf.String(), "- level_100")
		assert.Contains(t, buf.String(), "- week_streak")
	})
}
