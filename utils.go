package dlock

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

var newUUID = func() string {
	return uuid.NewV4().String()
}

var getCurrentTime = func() time.Time {
	return time.Now().UTC()
}
