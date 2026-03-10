package clierr

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const envDebug = "P_DEBUG"

type userMessager interface {
	UserMessage() string
}

// Wrap annotates an error with a user-facing message.
func Wrap(userMessage string, err error) error {
	if err == nil {
		return errors.New(userMessage)
	}
	return friendlyError{user: userMessage, err: err}
}

// Format returns a user-friendly error string. When P_DEBUG is enabled,
// it also includes the unwrap chain.
func Format(err error) string {
	if err == nil {
		return ""
	}

	user := err.Error()
	var um userMessager
	if errors.As(err, &um) {
		user = um.UserMessage()
	}

	if !debugEnabled() {
		return user
	}
	return user + "\n\nDebug:\n" + unwrapChain(err)
}

type friendlyError struct {
	user string
	err  error
}

func (e friendlyError) Error() string {
	return fmt.Sprintf("%s: %v", e.user, e.err)
}

func (e friendlyError) UserMessage() string {
	return e.user
}

func (e friendlyError) Unwrap() error {
	return e.err
}

func debugEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(envDebug)))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func unwrapChain(err error) string {
	lines := []string{}
	seen := map[error]bool{}
	cur := err
	for cur != nil {
		if seen[cur] {
			lines = append(lines, "(cycle detected)")
			break
		}
		seen[cur] = true
		lines = append(lines, "- "+cur.Error())
		cur = errors.Unwrap(cur)
	}
	return strings.Join(lines, "\n")
}
