package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteBadgeCommand(t *testing.T) {
	t.Run("Does not have badge", func(t *testing.T) {
		mockDB := new(MockDBClient)
		mockDB.On("GetUserBadges", mock.Anything, "user123").Return([]string{"other_badge"}, "John Doe", nil)
		SetDBClient(mockDB)
		defer SetDBClient(nil)
		userIDFlag = ""
		badgeFlag = ""
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"delete-badge", "user123", "week_streak"})
		err := RootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "does not have the badge")
		mockDB.AssertExpectations(t)
	})
	t.Run("Success path", func(t *testing.T) {
		mockDB := new(MockDBClient)
		mockDB.On("GetUserBadges", mock.Anything, "user123").Return([]string{"week_streak"}, "John Doe", nil)
		mockDB.On("DeleteBadgeFromUser", mock.Anything, "user123", "week_streak").Return(nil)
		SetDBClient(mockDB)
		defer SetDBClient(nil)
		userIDFlag = ""
		badgeFlag = ""
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"delete-badge", "user123", "week_streak"})
		err := RootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Success: Badge 'week_streak' removed from user 'John Doe'")
		mockDB.AssertExpectations(t)
	})
}
