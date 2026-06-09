package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListBadgesCommand(t *testing.T) {
	t.Run("No badges assigned", func(t *testing.T) {
		mockDB := new(MockDBClient)
		mockDB.On("GetUserBadges", mock.Anything, "user123").Return([]string{}, "John Doe", nil)
		SetDBClient(mockDB)
		defer SetDBClient(nil)
		userIDFlag = ""
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"list-badges", "user123"})
		err := RootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "User: John Doe (user123)")
		assert.Contains(t, buf.String(), "Badges: [No badges assigned]")
		mockDB.AssertExpectations(t)
	})
	t.Run("Multiple badges", func(t *testing.T) {
		mockDB := new(MockDBClient)
		mockDB.On("GetUserBadges", mock.Anything, "user123").Return([]string{"week_streak", "level_100"}, "John Doe", nil)
		SetDBClient(mockDB)
		defer SetDBClient(nil)
		userIDFlag = ""
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"list-badges", "user123"})
		err := RootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "User: John Doe (user123)")
		assert.Contains(t, buf.String(), "Badges (total 2):")
		assert.Contains(t, buf.String(), "- week_streak")
		assert.Contains(t, buf.String(), "- level_100")
		mockDB.AssertExpectations(t)
	})
	t.Run("DB error propagates", func(t *testing.T) {
		mockDB := new(MockDBClient)
		mockDB.On("GetUserBadges", mock.Anything, "user123").Return([]string{}, "", errors.New("connection failed"))
		SetDBClient(mockDB)
		defer SetDBClient(nil)
		userIDFlag = ""
		RootCmd.SetArgs([]string{"list-badges", "user123"})
		err := RootCmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection failed")
	})
}
