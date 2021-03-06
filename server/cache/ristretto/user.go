package ristretto

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/mazrean/Quantainer/cache"
	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/service"
)

type User struct {
	meCache     *ristretto.Cache
	activeUsers *ristretto.Cache
}

const (
	activeUsersKey = "active_users"
	activeUsersTTL = time.Hour
)

func NewUser() (*User, error) {
	meCache, err := ristretto.NewCache(&ristretto.Config{
		/*
			アクセス頻度を保持する要素の数。
			一般的には最大で格納される要素数の10倍程度が良いらしいが、
			最大でtraP部員数しか格納されないことを考えて500を設定する。
		*/
		NumCounters: 500,
		/*
			キャッシュの最大サイズ。
			あまり大きくしすぎるとメモリが足りなくなるので注意!
			*UserInfo1つあたり8Byteなので、8*500=20kB<2**15に設定する。
		*/
		MaxCost:     1 << 15,
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create meCache: %v", err)
	}

	activeUsers, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10,
		MaxCost:     64,
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create activeUsers: %v", err)
	}

	return &User{
		meCache:     meCache,
		activeUsers: activeUsers,
	}, nil
}

func (u *User) GetMe(ctx context.Context, accessToken values.OIDCAccessToken) (*service.UserInfo, error) {
	iUser, ok := u.meCache.Get(string(accessToken))
	if !ok {
		return nil, cache.ErrCacheMiss
	}

	user, ok := iUser.(*service.UserInfo)
	if !ok {
		return nil, fmt.Errorf("failed to cast meCache: %v", iUser)
	}

	return user, nil
}

func (u *User) SetMe(ctx context.Context, session *domain.OIDCSession, user *service.UserInfo) error {
	// キャッシュ追加待ちのキューに入るだけで、すぐにはキャッシュが効かないのに注意
	ok := u.meCache.SetWithTTL(
		string(session.GetAccessToken()),
		user,
		8,
		// sessionの有効期限が切れるとキャッシュが消えるようにTTLを設定する
		time.Until(session.GetExpiresAt()),
	)
	if !ok {
		return errors.New("failed to set meCache")
	}

	return nil
}

func (u *User) GetAllActiveUsers(ctx context.Context) ([]*service.UserInfo, error) {
	iUsers, ok := u.activeUsers.Get(activeUsersKey)
	if !ok {
		return nil, cache.ErrCacheMiss
	}

	users, ok := iUsers.([]*service.UserInfo)
	if !ok {
		return nil, fmt.Errorf("failed to cast activeUsers: %v", iUsers)
	}

	return users, nil
}

func (u *User) SetAllActiveUsers(ctx context.Context, users []*service.UserInfo) error {
	// キャッシュ追加待ちのキューに入るだけで、すぐにはキャッシュが効かないのに注意
	ok := u.activeUsers.SetWithTTL(
		activeUsersKey,
		users,
		1,
		activeUsersTTL,
	)
	if !ok {
		return errors.New("failed to set activeUsers")
	}

	return nil
}
