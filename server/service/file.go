package service

import (
	"context"
	"io"
)

type File interface {
	Upload(ctx context.Context, reader io.Reader) error
	Download(ctx context.Context, writer io.Writer) error
	/*
		GarbageCollection
		定時間resourceが作られなかったファイルを削除する
		基本的にサービス実装内のworkerが自動実行する
	*/
	GarbageCollection() error
}
