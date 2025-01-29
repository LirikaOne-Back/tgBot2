package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
	"tgBot/lib/e"
	"tgBot/storage"
)

const (
	RndCmd       = "/rnd"
	Help         = "/help"
	StartCmd     = "/start"
	RemoveRndCmd = "/remove"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case RemoveRndCmd:
		return p.remove(chatID, username)
	case Help:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err = p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err = p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) remove(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: remove", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSaved) {
		return err
	}
	if errors.Is(err, storage.ErrNoSaved) {
		return p.tg.SendMessage(chatID, msgNoRemove)
	}

	if err = p.tg.SendMessage(chatID, msgRemove); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSaved) {
		return err
	}
	if errors.Is(err, storage.ErrNoSaved) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err = p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
