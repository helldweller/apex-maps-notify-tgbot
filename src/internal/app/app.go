package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/sync/errgroup"

	"package/main/internal/apexapi"
)

var maps apexapi.Maps
var ctx, cancel = context.WithCancel(context.Background())
var group, groupCtx = errgroup.WithContext(ctx)

func Run() {

	log.Info("Starting app")

	group.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		defer close(signalChannel)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
		select {
		case sig := <-signalChannel:
			log.Errorf("Received signal: %s", sig)
			cancel()
		case <-groupCtx.Done():
			log.Error("Closing signal goroutine")
			return groupCtx.Err()
		}
		return nil
	})

	group.Go(func() error {
		if err := maps.Update(conf.ApexApiKey); err != nil {
			log.Error(err) // shit
		} // get on start and then every interval seconds
		interval := 60 // sec
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case <-ticker.C:
				if err := maps.Update(conf.ApexApiKey); err != nil {
					log.Error(err)
				}
			case <-groupCtx.Done():
				log.Error("Closing apexmaps goroutine")
				return groupCtx.Err()
			}
		}
		return nil
	})

	bot, err := tgbotapi.NewBotAPI(conf.BotApiKey)
	if err != nil {
		log.Panic(err)
	}
	tgbotapi.SetLogger(log)
	bot.Debug = conf.BotDebug
	log.Infof("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	group.Go(func() error {
		for {
			select {
			case update := <-updates:
				if update.Message == nil {
					continue
				}
				if !update.Message.IsCommand() {
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

				switch update.Message.Command() {
				case "help":
					msg.Text = "I understand /map"
				case "map":
					now := time.Now()
					nextStartAt := time.Unix(maps.Next.Start, 0)
					nextEndAt := time.Unix(maps.Next.End, 0)
					nextDiff := nextStartAt.Sub(now)
					nextLasts := nextEndAt.Sub(nextStartAt)
					msg.Text = fmt.Sprintf("Карта сейчас *%s*\nСледующая карта *%s* через *%dч %dм* и продлится *%dч %dм*",
						maps.Current.Map,
						maps.Next.Map,
						int(nextDiff.Hours()),
						int(nextDiff.Minutes())-int(nextDiff.Hours())*60,
						int(nextLasts.Hours()),
						int(nextLasts.Minutes())-int(nextLasts.Hours())*60,
					)
					msg.ReplyToMessageID = update.Message.MessageID
					msg.ParseMode = "MarkdownV2"
					log.Infof("Recived new message from user %s, chat ID %d", update.Message.From, update.Message.Chat.ID)
				default:
					msg.Text = "I don't know that command, use /help"
				}
				if _, err := bot.Send(msg); err != nil {
					log.Error(err)
				}
			case <-groupCtx.Done():
				log.Error("Closing tgBot goroutine")
				return groupCtx.Err()
			}
		}
		return nil
	})

	err = group.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Error("Context was canceled")
		} else {
			log.Errorf("Received error: %v\n", err)
		}
	} else {
		log.Error("Sucsessfull finished")
	}

}
