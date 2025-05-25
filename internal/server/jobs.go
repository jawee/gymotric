package server

import (
	"context"
	"log/slog"
	"time"
	"weight-tracker/internal/utils"

	_ "github.com/joho/godotenv/autoload"
)

func (s *Server) RegisterJobs() {
	go s.cleanupUnverifiedUsers()
}

func (s *Server) cleanupUnverifiedUsers() {
	first := true
	for {
		if !first {
			time.Sleep(utils.AccountConfirmationTokenExpireMinutes * time.Minute)
		}
		slog.Info("Running cleanup of unverified users")

		context := context.Background()
		users, err := s.db.GetRepository().GetUnverifiedUsers(context)

		if err != nil {
			slog.Error("Failed to get unverified users", "error", err)
			continue
		}
		slog.Info("Found unverified users", "count", len(users))
		for _, user := range users {
			createdTime, err := time.Parse(time.RFC3339, user.CreatedOn)

			if err != nil {
				slog.Error("Failed to parse user creation time", "userID", user.ID, "error", err)
				continue
			}
			if time.Since(createdTime) < utils.AccountConfirmationTokenExpireMinutes*time.Minute {
				slog.Info("Skipping deletion for user", "userID", user.ID, "createdOn", createdTime)
				continue
			}

			slog.Info("Deleting unverified user", "userID", user.ID)
			_, err = s.db.GetRepository().DeleteUser(context, user.ID)
			if err != nil {
				slog.Error("Failed to delete unverified user", "userID", user.ID, "error", err)
				continue
			}
			slog.Info("Deleted unverified user", "userID", user.ID)
		}

		slog.Info("Cleanup of unverified users completed successfully")
		first = false
	}
}
