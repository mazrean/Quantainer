package ristretto

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mazrean/Quantainer/cache"
	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/service"
	"github.com/stretchr/testify/assert"
)

func TestGetMe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	userCache, err := NewUser()
	if err != nil {
		t.Fatalf("failed to create user cache: %v", err)
	}

	type test struct {
		description string
		keyExist    bool
		valueBroken bool
		userInfo    *service.UserInfo
		isErr       bool
		err         error
	}

	testCases := []test{
		{
			description: "特に問題ないのでエラーなし",
			keyExist:    true,
			userInfo: service.NewUserInfo(
				values.NewTrapMemberID(uuid.New()),
				values.NewTrapMemberName("mazrean"),
				values.TrapMemberStatusActive,
			),
		},
		{
			description: "キーが存在しないのでErrCacheMiss",
			keyExist:    false,
			isErr:       true,
			err:         cache.ErrCacheMiss,
		},
		{
			// 実際には発生しないが念の為確認
			description: "値が壊れているのでエラー",
			keyExist:    true,
			valueBroken: true,
			isErr:       true,
		},
	}

	for i, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			accessToken := values.NewOIDCAccessToken(fmt.Sprintf("access token%d", i))
			if testCase.keyExist {
				if testCase.valueBroken {
					ok := userCache.meCache.Set(string(accessToken), "broken", 8)
					assert.True(t, ok)

					userCache.meCache.Wait()
				} else {
					ok := userCache.meCache.Set(string(accessToken), testCase.userInfo, 8)
					assert.True(t, ok)

					userCache.meCache.Wait()
				}
			}

			user, err := userCache.GetMe(ctx, accessToken)

			if testCase.isErr {
				if testCase.err == nil {
					assert.Error(t, err)
				} else if !errors.Is(err, testCase.err) {
					t.Errorf("error must be %v, but actual is %v", testCase.err, err)
				}
			} else {
				assert.NoError(t, err)
			}
			if err != nil {
				return
			}

			assert.Equal(t, testCase.userInfo, user)
		})
	}
}

func TestSetMe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	userCache, err := NewUser()
	if err != nil {
		t.Fatalf("failed to create user cache: %v", err)
	}

	type test struct {
		description string
		keyExist    bool
		beforeValue *service.UserInfo
		session     *domain.OIDCSession
		userInfo    *service.UserInfo
		ttl         time.Duration
		isErr       bool
		err         error
	}

	now := time.Now()

	testCases := []test{
		{
			description: "特に問題ないのでエラーなし",
			session: domain.NewOIDCSession(
				values.NewOIDCAccessToken("access token1"),
				now.Add(2*time.Second),
			),
			userInfo: service.NewUserInfo(
				values.NewTrapMemberID(uuid.New()),
				values.NewTrapMemberName("mazrean"),
				values.TrapMemberStatusActive,
			),
			ttl: 2 * time.Second,
		},
		{
			description: "元からキーがあっても上書きする",
			keyExist:    true,
			beforeValue: service.NewUserInfo(
				values.NewTrapMemberID(uuid.New()),
				values.NewTrapMemberName("mazrean"),
				values.TrapMemberStatusActive,
			),
			session: domain.NewOIDCSession(
				values.NewOIDCAccessToken("access token2"),
				now.Add(2*time.Second),
			),
			userInfo: service.NewUserInfo(
				values.NewTrapMemberID(uuid.New()),
				values.NewTrapMemberName("mazrean"),
				values.TrapMemberStatusActive,
			),
			ttl: 2 * time.Second,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			if testCase.keyExist {
				ok := userCache.meCache.Set(string(testCase.session.GetAccessToken()), testCase.beforeValue, 8)
				assert.True(t, ok)

				userCache.meCache.Wait()
			}

			err := userCache.SetMe(ctx, testCase.session, testCase.userInfo)

			if testCase.isErr {
				if testCase.err == nil {
					assert.Error(t, err)
				} else if !errors.Is(err, testCase.err) {
					t.Errorf("error must be %v, but actual is %v", testCase.err, err)
				}
			} else {
				assert.NoError(t, err)
			}
			if err != nil {
				return
			}

			// キャッシュが設定されるまで待機
			userCache.meCache.Wait()

			// OIDCSessionの期限前なのでキャッシュされている
			value, ok := userCache.meCache.Get(string(testCase.session.GetAccessToken()))
			assert.True(t, ok)
			assert.Equal(t, testCase.userInfo, value)

			<-time.NewTimer(testCase.ttl).C

			// OIDCSessionの期限が切れたらキャッシュは削除される
			_, ok = userCache.meCache.Get(string(testCase.session.GetAccessToken()))
			assert.False(t, ok)
		})
	}
}

func TestGetAllActiveUsers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	userCache, err := NewUser()
	if err != nil {
		t.Fatalf("failed to create user cache: %v", err)
	}

	type test struct {
		description string
		keyExist    bool
		valueBroken bool
		users       []*service.UserInfo
		isErr       bool
		err         error
	}

	testCases := []test{
		{
			description: "特に問題ないのでエラーなし",
			keyExist:    true,
			users: []*service.UserInfo{
				service.NewUserInfo(
					values.NewTrapMemberID(uuid.New()),
					values.NewTrapMemberName("mazrean"),
					values.TrapMemberStatusActive,
				),
			},
		},
		{
			description: "ユーザー数が500人でも問題なし",
			keyExist:    true,
			users:       make([]*service.UserInfo, 500),
		},
		{
			description: "キーが存在しないのでErrCacheMiss",
			keyExist:    false,
			isErr:       true,
			err:         cache.ErrCacheMiss,
		},
		{
			// 実際には発生しないが念の為確認
			description: "値が壊れているのでエラー",
			keyExist:    true,
			valueBroken: true,
			isErr:       true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			if testCase.keyExist {
				if testCase.valueBroken {
					ok := userCache.activeUsers.Set(activeUsersKey, "broken", 8)
					assert.True(t, ok)
				} else {
					ok := userCache.activeUsers.Set(activeUsersKey, testCase.users, 8)
					assert.True(t, ok)
				}

				userCache.activeUsers.Wait()
				defer userCache.activeUsers.Clear()
			}

			users, err := userCache.GetAllActiveUsers(ctx)

			if testCase.isErr {
				if testCase.err == nil {
					assert.Error(t, err)
				} else if !errors.Is(err, testCase.err) {
					t.Errorf("error must be %v, but actual is %v", testCase.err, err)
				}
			} else {
				assert.NoError(t, err)
			}
			if err != nil {
				return
			}

			assert.Equal(t, testCase.users, users)
		})
	}
}

func TestSetAllActiveUsers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	userCache, err := NewUser()
	if err != nil {
		t.Fatalf("failed to create user cache: %v", err)
	}

	type test struct {
		description string
		keyExist    bool
		beforeValue []*service.UserInfo
		users       []*service.UserInfo
		isErr       bool
		err         error
	}

	testCases := []test{
		{
			description: "特に問題ないのでエラーなし",
			users: []*service.UserInfo{
				service.NewUserInfo(
					values.NewTrapMemberID(uuid.New()),
					values.NewTrapMemberName("mazrean"),
					values.TrapMemberStatusActive,
				),
			},
		},
		{
			description: "ユーザー数が500人でもエラーなし",
			users:       make([]*service.UserInfo, 500),
		},
		{
			description: "元からキーがあっても上書きする",
			keyExist:    true,
			beforeValue: []*service.UserInfo{
				service.NewUserInfo(
					values.NewTrapMemberID(uuid.New()),
					values.NewTrapMemberName("mazrean"),
					values.TrapMemberStatusActive,
				),
			},
			users: []*service.UserInfo{
				service.NewUserInfo(
					values.NewTrapMemberID(uuid.New()),
					values.NewTrapMemberName("mazrean"),
					values.TrapMemberStatusActive,
				),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			if testCase.keyExist {
				ok := userCache.activeUsers.Set(activeUsersKey, testCase.beforeValue, 1)
				assert.True(t, ok)

				userCache.activeUsers.Wait()
			}

			err := userCache.SetAllActiveUsers(ctx, testCase.users)

			if testCase.isErr {
				if testCase.err == nil {
					assert.Error(t, err)
				} else if !errors.Is(err, testCase.err) {
					t.Errorf("error must be %v, but actual is %v", testCase.err, err)
				}
			} else {
				assert.NoError(t, err)
			}
			if err != nil {
				return
			}

			// キャッシュが設定されるまで待機
			userCache.activeUsers.Wait()

			// OIDCSessionの期限前なのでキャッシュされている
			value, ok := userCache.activeUsers.Get(activeUsersKey)
			assert.True(t, ok)
			assert.IsType(t, []*service.UserInfo{}, value)
			actualUsers := value.([]*service.UserInfo)

			for i, user := range testCase.users {
				if user == nil {
					assert.Nil(t, actualUsers[i])
				} else {
					assert.Equal(t, *actualUsers[i], *user)
				}
			}
		})
	}
}
