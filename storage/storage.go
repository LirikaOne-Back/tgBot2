package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"tgBot/lib/e"
	"time"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
	GiveAll(ctx context.Context, userName string) (*[]Page, error)
}

var ErrNoSaved = errors.New("no saved page")

type Page struct {
	URL      string
	UserName string
	Time     time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
