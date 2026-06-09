package cmd

import (
	"context"
	"os"

	"github.com/stretchr/testify/mock"
)

func init() {
	_ = os.Setenv("DATABASE_URL", "mongodb://localhost:27017/test")
	_ = os.Setenv("JORBITES_URL", "http://localhost:3000")
}

type MockDBClient struct {
	mock.Mock
}

func (m *MockDBClient) AssignBadgeToUser(ctx context.Context, userID string, badgeName string) error {
	args := m.Called(ctx, userID, badgeName)
	return args.Error(0)
}
func (m *MockDBClient) DeleteBadgeFromUser(ctx context.Context, userID string, badgeName string) error {
	args := m.Called(ctx, userID, badgeName)
	return args.Error(0)
}
func (m *MockDBClient) GetUserBadges(ctx context.Context, userID string) ([]string, string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.String(1), args.Error(2)
}
func (m *MockDBClient) Close(ctx context.Context) error {
	return nil
}
