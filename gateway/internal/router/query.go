package router

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
)

func parseLimit(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func boundedLimit(raw string, fallback int) int {
	parsed := parseLimit(raw, fallback)
	if parsed > maxReplayLimit {
		return maxReplayLimit
	}
	return parsed
}

func parseInt64(raw string, fallback int64) int64 {
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fallback
	}
	return v
}

func validateEventQuery(query EventQuery) error {
	if query.AgentID != "" {
		if _, err := uuid.Parse(query.AgentID); err != nil {
			return err
		}
	}
	if query.SinceMs > 0 && query.UntilMs > 0 && query.SinceMs > query.UntilMs {
		return errors.New("sinceMs must be <= untilMs")
	}
	return nil
}
