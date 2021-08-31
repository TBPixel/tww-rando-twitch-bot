package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/timshannon/badgerhold"
)

type UserQueryField string

const (
	FieldID              UserQueryField = "ID"
	FieldTwitchID        UserQueryField = "TwitchID"
	FieldRacetimeID      UserQueryField = "RacetimeID"
	FieldActiveInChannel UserQueryField = "ActiveInChannel"
)

var (
	ErrNotFound = errors.New("no results match that query")
	ErrExists   = errors.New("resource already exists")
)

// User
type User struct {
	ID                uint64 `badgerhold:"key"`
	TwitchID          string `badgerhold:"unique"`
	RacetimeID        string
	TwitchName        string
	TwitchDisplayName string
	ProfileImageURL   string
	ActiveInChannel   bool
	JoinedAt          time.Time
}

type UserQuery struct {
	Field UserQueryField
	Value interface{}
}

type UserUpdate struct {
	TwitchID          *string
	RacetimeID        *string
	TwitchName        *string
	TwitchDisplayName *string
	ActiveInChannel   *bool
}

// FindUser
func (db *DB) FindUsers(query UserQuery) ([]*User, error) {
	var q *badgerhold.Query
	if query.Field == FieldID {
		q = badgerhold.Where(badgerhold.Key).Eq(query.Value).Limit(1)
	} else {
		q = badgerhold.Where(string(query.Field)).Eq(query.Value)
	}

	var users []*User
	err := db.store.Find(&users, q)
	if err != nil {
		return nil, fmt.Errorf("error while looking up users %w with query %+v\n", err, query)
	}

	return users, nil
}

func (db *DB) FindUser(query UserQuery) (*User, error) {
	users, err := db.FindUsers(query)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	return users[0], nil
}

// CreateUser inserts a new user into the database
func (db *DB) CreateUser(twitchID, twitchName, twitchDisplayName, profileImageURL string) (*User, error) {
	users, err := db.FindUsers(UserQuery{
		Field: FieldTwitchID,
		Value: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("error while creating user: %w", err)
	}

	if users != nil && len(users) > 0 {
		return nil, ErrExists
	}

	user := User{
		TwitchID:          twitchID,
		TwitchName:        twitchName,
		TwitchDisplayName: twitchDisplayName,
		ProfileImageURL:   profileImageURL,
		ActiveInChannel:   false,
		JoinedAt:          time.Now(),
	}
	err = db.store.Insert(badgerhold.NextSequence(), &user)
	if err != nil {
		return nil, fmt.Errorf("error while creating new user: %w", err)
	}

	return &user, nil
}

// UpdateUser
func (db *DB) UpdateUser(id uint64, user UserUpdate) (*User, error) {
	users, err := db.FindUsers(UserQuery{
		Field: FieldID,
		Value: id,
	})
	if err != nil {
		return nil, fmt.Errorf("error while updating user: %w", err)
	}
	if len(users) == 0 {
		return nil, ErrNotFound
	}

	u := users[0]
	if user.TwitchID != nil {
		u.TwitchID = *user.TwitchID
	}
	if user.TwitchName != nil {
		u.TwitchName = *user.TwitchName
	}
	if user.TwitchDisplayName != nil {
		u.TwitchDisplayName = *user.TwitchDisplayName
	}
	if user.RacetimeID != nil {
		u.RacetimeID = *user.RacetimeID
	}
	if user.ActiveInChannel != nil {
		u.ActiveInChannel = *user.ActiveInChannel
	}

	err = db.store.Update(id, u)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return u, nil
}

// DeleteUser
func (db *DB) DeleteUser(id uint64) error {
	err := db.store.Delete(id, &User{})
	if err != nil {
		return fmt.Errorf("error deleting user with id %v", id)
	}

	return nil
}
