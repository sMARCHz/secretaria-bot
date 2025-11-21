package services

import "github.com/sMARCHz/go-secretaria-bot/internal/core/errors"

// CommandHandler handles a specific command namespace (e.g. finance).
type CommandHandler interface {
    // Match returns true when this handler should process the command token.
    Match(cmd string) bool

    // Handle executes the command. msgArgs is tokenized input (fields).
    Handle(msgArgs []string) (string, *errors.AppError)
}
