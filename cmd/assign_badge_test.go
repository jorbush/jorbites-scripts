package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAssignBadgeCommand(t *testing.T) {
	t.Run("Already has badge", func(t *testing.T) {
		mockDB := new(MockDBClient)
		mockDB.On("GetUserBadges", mock.Anything, "user123").Return([]string{"week_streak"}, "John Doe", nil)
		SetDBClient(mockDB)
		defer SetDBClient(nil)
		// Set flags/args
		userIDFlag = ""
		badgeFlag = ""
		forceFlag = false
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"assign-badge", "user123", "week_streak"})
		err := RootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "already has the badge")
		mockDB.AssertExpectations(t)
	})
	t.Run("Success path with API validation", func(t *testing.T) {
		// Start mock HTTP server for validation
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/badges", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]string{"week_streak", "level_100"})
		}))
		defer server.Close()
		SetAppURL(server.URL)
		defer SetAppURL("http://localhost:3000")
		mockDB := new(MockDBClient)
		mockDB.On("GetUserBadges", mock.Anything, "user123").Return([]string{}, "John Doe", nil)
		mockDB.On("AssignBadgeToUser", mock.Anything, "user123", "week_streak").Return(nil)
		SetDBClient(mockDB)
		defer SetDBClient(nil)
		userIDFlag = ""
		badgeFlag = ""
		forceFlag = false
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"assign-badge", "user123", "week_streak"})
		err := RootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Success: Badge 'week_streak' assigned to user 'John Doe'")
		mockDB.AssertExpectations(t)
	})
	t.Run("API validation fails but user confirms", func(t *testing.T) {
		// Mock API returns other badges, not 'special_badge'
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]string{"week_streak"})
		}))
		defer server.Close()
		SetAppURL(server.URL)
		defer SetAppURL("http://localhost:3000")
		// Mock user inputs 'y'
		readerScanner = bufio.NewScanner(strings.NewReader("y\n"))
		defer func() { readerScanner = nil }()
		mockDB := new(MockDBClient)
		mockDB.On("GetUserBadges", mock.Anything, "user123").Return([]string{}, "John Doe", nil)
		mockDB.On("AssignBadgeToUser", mock.Anything, "user123", "special_badge").Return(nil)
		SetDBClient(mockDB)
		defer SetDBClient(nil)
		userIDFlag = ""
		badgeFlag = ""
		forceFlag = false
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"assign-badge", "user123", "special_badge"})
		err := RootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Success: Badge 'special_badge' assigned to user 'John Doe'")
		mockDB.AssertExpectations(t)
	})
	t.Run("API validation fails and user aborts", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]string{"week_streak"})
		}))
		defer server.Close()
		SetAppURL(server.URL)
		defer SetAppURL("http://localhost:3000")
		// Mock user inputs 'n'
		readerScanner = bufio.NewScanner(strings.NewReader("n\n"))
		defer func() { readerScanner = nil }()
		mockDB := new(MockDBClient)
		mockDB.On("GetUserBadges", mock.Anything, "user123").Return([]string{}, "John Doe", nil)
		SetDBClient(mockDB)
		defer SetDBClient(nil)
		userIDFlag = ""
		badgeFlag = ""
		forceFlag = false
		RootCmd.SetArgs([]string{"assign-badge", "user123", "special_badge"})
		err := RootCmd.Execute()
		assert.Error(t, err)
		assert.Equal(t, "assignment aborted by user", err.Error())
	})
}
