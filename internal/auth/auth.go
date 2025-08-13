package auth

import (
	"context"
	"errors"
)

// ctxKey is an unexported type to prevent context key collisions.
type ctxKey string

const identityKey = ctxKey("identity")

// ErrUnauthorized is a standard error for when a user is not found in the context.
var ErrUnauthorized = errors.New("unauthorized: user identity not found")

// Identity represents the authenticated user. It's a struct for future extensibility.
type Identity struct {
	UserID string
}

// ToContext adds an Identity to the given context.
func ToContext(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, identityKey, id)
}

// FromContext retrieves an Identity from the context.
func FromContext(ctx context.Context) (Identity, error) {
	identity, ok := ctx.Value(identityKey).(Identity)
	if !ok {
		return Identity{}, ErrUnauthorized
	}
	return identity, nil
}
