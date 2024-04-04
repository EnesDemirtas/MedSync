// Package mid contains the set of values the middleware is responsible
// to extract and set.
package mid

import (
	"context"
	"errors"

	"github.com/EnesDemirtas/medisync/business/api/auth"
	"github.com/EnesDemirtas/medisync/business/core/crud/inventorybus"
	"github.com/EnesDemirtas/medisync/business/core/crud/medicinebus"
	"github.com/EnesDemirtas/medisync/business/core/crud/tagbus"
	"github.com/EnesDemirtas/medisync/business/core/crud/userbus"
	"github.com/google/uuid"
)

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
	userKey
	tagKey
	medicineKey
	inventoryKey
)

func SetClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

func GetClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

// GetUserID returns the claims from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}

	return v, nil
}

// GetUser returns the user from the context.
func GetUser(ctx context.Context) (userbus.User, error) {
	v, ok := ctx.Value(userKey).(userbus.User)
	if !ok {
		return userbus.User{}, errors.New("user not found in context")
	}

	return v, nil
}

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func SetUser(ctx context.Context, usr userbus.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

// GetTag returns the tag from the context.
func GetTag(ctx context.Context) (tagbus.Tag, error) {
	v, ok := ctx.Value(tagKey).(tagbus.Tag)
	if !ok {
		return tagbus.Tag{}, errors.New("tag not found in context")
	}

	return v, nil
}

func SetTag(ctx context.Context, tag tagbus.Tag) context.Context {
	return context.WithValue(ctx, tagKey, tag)
}

// GetMedicine returns the medicine from the context.
func GetMedicine(ctx context.Context) (medicinebus.Medicine, error) {
	v, ok := ctx.Value(medicineKey).(medicinebus.Medicine)
	if !ok {
		return medicinebus.Medicine{}, errors.New("medicine not found in context")
	}

	return v, nil
}

func SetMedicine(ctx context.Context, med medicinebus.Medicine) context.Context {
	return context.WithValue(ctx, medicineKey, med)
}

// GetInventory returns the inventory from the context.
func GetInventory(ctx context.Context) (inventorybus.Inventory, error) {
	v, ok := ctx.Value(inventoryKey).(inventorybus.Inventory)
	if !ok {
		return inventorybus.Inventory{}, errors.New("inventory not found in context")
	}

	return v, nil
}

func SetInventory(ctx context.Context, inv inventorybus.Inventory) context.Context {
	return context.WithValue(ctx, inventoryKey, inv)
}