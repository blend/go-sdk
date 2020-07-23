package web

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLocalSessionCache(t *testing.T) {
	assert := assert.New(t)

	lsc := NewLocalSessionCache()

	session := &Session{UserID: "bailey", SessionID: NewSessionID()}
	assert.Nil(lsc.PersistHandler(context.TODO(), session))

	fetched, err := lsc.FetchHandler(context.TODO(), session.SessionID)
	assert.Nil(err)
	assert.Equal(session.UserID, fetched.UserID)

	assert.Nil(lsc.RemoveHandler(context.TODO(), session.SessionID))

	removed, err := lsc.FetchHandler(context.TODO(), session.SessionID)
	assert.Nil(err)
	assert.Nil(removed)
}
