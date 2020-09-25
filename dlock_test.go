package dlock

import (
	"context"
	"dlock/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"time"
)

func Test_Take(t *testing.T) {
	controller := gomock.NewController(t)
	mockStore := mocks.NewMockLockStore(controller)
	ctx := context.Background()
	oldNewUUID := newUUID
	oldgetCurrentTime := getCurrentTime
	currTime := getCurrentTime()
	uid := newUUID()
	defer func() {
		newUUID = oldNewUUID
		getCurrentTime = oldgetCurrentTime
		controller.Finish()
	}()
	newUUID = func() string {
		return uid
	}
	getCurrentTime = func() time.Time {
		return currTime
	}
	type args struct {
		key string
		ttl int64
	}
	type test struct {
		name string
		args
		mock        func()
		lock        *Lock
		wantErr     bool
		wantLocked  bool
		expectedErr error
	}
	tests := []test{
		{
			name: "take lock success",
			args: args{
				key: "key",
				ttl: 10,
			},
			lock: &Lock{
				ctx:   ctx,
				store: mockStore,
				uid:   uid,
			},
			mock: func() {
				mockStore.EXPECT().Set(ctx, "key", uid, getCurrentTime().Add(time.Duration(10)*time.Second)).Return(nil)
			},
			wantLocked: true,
		},
		{
			name: "take lock error lock already set",
			args: args{
				key: "key",
				ttl: 1000,
			},
			lock: &Lock{
				ctx:    ctx,
				store:  mockStore,
				uid:    uid,
				locked: true,
				key:    "key",
				expiry: currTime.Add(10 * time.Second),
			},
			mock: func() {
				mockStore.EXPECT().Get(ctx, "key").Return(uid)
			},
			wantErr:     true,
			expectedErr: errors.New(lockAlreadySet),
			wantLocked:  true,
		},
		{
			name: "take lock error in setting lock",
			args: args{
				key: "key",
				ttl: 1000,
			},
			lock: &Lock{
				ctx:   ctx,
				store: mockStore,
				uid:   uid,
			},
			mock: func() {
				mockStore.EXPECT().Set(ctx, "key", uid, currTime.Add(time.Second*1000)).Return(errors.New("error"))
			},
			wantErr:     true,
			expectedErr: errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.lock.Take(tt.key, tt.ttl)
			if tt.wantErr && !reflect.DeepEqual(tt.expectedErr, err) {
				t.Errorf("got:  %#v\n\t want: %#v", err, tt.expectedErr)
			}
			if tt.wantLocked != tt.lock.locked {
				t.Errorf("got:  %#v\n\t want: %#v", tt.lock.locked, tt.wantLocked)
			}
		})
	}
}

func Test_Release(t *testing.T) {
	controller := gomock.NewController(t)
	mockStore := mocks.NewMockLockStore(controller)
	ctx := context.Background()
	oldNewUUID := newUUID
	oldgetCurrentTime := getCurrentTime
	currTime := getCurrentTime()
	uid := newUUID()
	defer func() {
		newUUID = oldNewUUID
		getCurrentTime = oldgetCurrentTime
		controller.Finish()
	}()
	newUUID = func() string {
		return uid
	}
	getCurrentTime = func() time.Time {
		return currTime
	}
	type test struct {
		name        string
		key         string
		mock        func()
		lock        *Lock
		wantErr     bool
		wantLocked  bool
		expectedErr error
	}
	tests := []test{
		{
			name: "release lock success",
			key:  "key",
			mock: func() {
				mockStore.EXPECT().Get(ctx, "key").Return(uid)
				mockStore.EXPECT().Delete(ctx, "key").Return(nil)
			},
			lock: &Lock{
				ctx:    ctx,
				key:    "key",
				uid:    uid,
				locked: true,
				expiry: currTime.Add(time.Second * 10),
				store:  mockStore,
			},
		},
		{
			name: "release lock error lock not set",
			key:  "key",
			mock: func() {},
			lock: &Lock{
				ctx:   ctx,
				store: mockStore,
				uid:   uid,
			},
			wantErr:     true,
			expectedErr: errors.New(lockNotSet),
		},
		{
			name: "release lock error in delete",
			key:  "key",
			mock: func() {
				mockStore.EXPECT().Get(ctx, "key").Return(uid)
				mockStore.EXPECT().Delete(ctx, "key").Return(errors.New("error"))
			},
			lock: &Lock{
				ctx:    ctx,
				key:    "key",
				uid:    uid,
				locked: true,
				expiry: currTime.Add(time.Second * 10),
				store:  mockStore,
			},
			wantLocked:  true,
			wantErr:     true,
			expectedErr: errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.lock.Release()
			if tt.wantErr && !reflect.DeepEqual(tt.expectedErr, err) {
				t.Errorf("got:  %#v\n\t want: %#v", err, tt.expectedErr)
			}
			if tt.wantLocked != tt.lock.locked {
				t.Errorf("got:  %#v\n\t want: %#v", tt.lock.locked, tt.wantLocked)
			}
		})
	}
}

func Test_isSet(t *testing.T) {
	controller := gomock.NewController(t)
	mockStore := mocks.NewMockLockStore(controller)
	ctx := context.Background()
	oldNewUUID := newUUID
	oldgetCurrentTime := getCurrentTime
	currTime := getCurrentTime()
	uid := newUUID()
	defer func() {
		newUUID = oldNewUUID
		getCurrentTime = oldgetCurrentTime
		controller.Finish()
	}()
	newUUID = func() string {
		return uid
	}
	getCurrentTime = func() time.Time {
		return currTime
	}
	type test struct {
		name        string
		lock        *Lock
		mock        func()
		wantLocked  bool
		expectedRes bool
	}
	tests := []test{
		{
			name: "lock set",
			lock: &Lock{
				ctx:    ctx,
				locked: true,
				key:    "key",
				uid:    uid,
				expiry: currTime.Add(time.Second * 10),
				store:  mockStore,
			},
			mock: func() {
				mockStore.EXPECT().Get(ctx, "key").Return(uid)
			},
			wantLocked:  true,
			expectedRes: true,
		},
		{
			name: "lock not set",
			lock: &Lock{},
			mock: func() {},
		},
		{
			name: "lock not set in store",
			lock: &Lock{
				ctx:    ctx,
				locked: true,
				key:    "key",
				uid:    uid,
				expiry: currTime.Add(time.Second * 10),
				store:  mockStore,
			},
			mock: func() {
				mockStore.EXPECT().Get(ctx, "key").Return("")
			},
		},
		{
			name: "lock expired",
			lock: &Lock{
				ctx:    ctx,
				locked: true,
				key:    "key",
				uid:    uid,
				expiry: currTime.Add(-10 * time.Second),
				store:  mockStore,
			},
			mock: func() {
				mockStore.EXPECT().Get(ctx, "key").Return(uid)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res := tt.lock.isSet()
			if tt.expectedRes != res {
				t.Errorf("got:  %#v\n\t want: %#v", res, tt.expectedRes)
			}
			if tt.wantLocked != tt.lock.locked {
				t.Errorf("got:  %#v\n\t want: %#v", tt.lock.locked, tt.wantLocked)
			}
		})
	}
}
