package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/sync/errgroup"

	"package/main/internal/apexapi"
)

var modes apexapi.Modes

var ctx, cancel = context.WithCancel(context.Background())
var group, groupCtx = errgroup.WithContext(ctx)

// Run func is similar to the main func.
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
		if err := modes.Update(conf.ApexAPIKey); err != nil {
			log.Error(err) // shit
		} // get on start and then every interval seconds
		interval := 120 // sec
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case <-ticker.C:
				if err := modes.Update(conf.ApexAPIKey); err != nil {
					log.Error(err)
				}
			case <-groupCtx.Done():
				log.Error("Closing apexmaps modes goroutine")
				return groupCtx.Err()
			}
		}
	})

	bot, err := tgbotapi.NewBotAPI(conf.BotAPIKey)
	md2regex := regexp.MustCompile(`(\_|\*|\[|\]|\(|\)|\~|\>|\#|\+|\-|\=|\||\{|\}|\.|\!)`)
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
					msg.Text = "I understand /map /ranked /ltm"
				case "map":
					now := time.Now()
					nextStartAt := time.Unix(modes.Pub.Next.Start, 0)
					nextEndAt := time.Unix(modes.Pub.Next.End, 0)
					nextDiff := nextStartAt.Sub(now)
					nextLasts := nextEndAt.Sub(nextStartAt)
					msg.Text = fmt.Sprintf("Карта сейчас *%s* и продлится *%dч %dм*\nСледующая карта *%s* и продлится *%dч %dм*\n[](%s)",
						md2regex.ReplaceAllString(modes.Pub.Current.Map, `\$1`),
						int(nextDiff.Hours()),
						int(nextDiff.Minutes())-int(nextDiff.Hours())*60,
					        md2regex.ReplaceAllString(modes.Pub.Next.Map, `\$1`),
						int(nextLasts.Hours()),
						int(nextLasts.Minutes())-int(nextLasts.Hours())*60,
						modes.Pub.Current.Asset,
					)
					msg.ReplyToMessageID = update.Message.MessageID
					msg.ParseMode = "MarkdownV2"
					log.Infof("Recived new message from user %s, chat ID %d", update.Message.From, update.Message.Chat.ID)
				case "ranked":
					now := time.Now()
					nextStartAt := time.Unix(modes.Ranked.Next.Start, 0)
					nextEndAt := time.Unix(modes.Ranked.Next.End, 0)
					nextDiff := nextStartAt.Sub(now)
					nextLasts := nextEndAt.Sub(nextStartAt)
					msg.Text = fmt.Sprintf("Карта в рейтинге сейчас *%s* и продлится *%dч %dм*\nСледующая карта *%s* и продлится *%dч %dм*",
						md2regex.ReplaceAllString(modes.Ranked.Current.Map, `\$1`),
						int(nextDiff.Hours()),
						int(nextDiff.Minutes())-int(nextDiff.Hours())*60,
						md2regex.ReplaceAllString(modes.Ranked.Next.Map, `\$1`),
						int(nextLasts.Hours()),
						int(nextLasts.Minutes())-int(nextLasts.Hours())*60,
					)
					msg.ReplyToMessageID = update.Message.MessageID
					msg.ParseMode = "MarkdownV2"
					log.Infof("Recived new message from user %s, chat ID %d", update.Message.From, update.Message.Chat.ID)
				case "ltm":
					now := time.Now()
					nextStartAt := time.Unix(modes.Ltm.Next.Start, 0)
					nextEndAt := time.Unix(modes.Ltm.Next.End, 0)
					nextDiff := nextStartAt.Sub(now)
					nextLasts := nextEndAt.Sub(nextStartAt)
					msg.Text = fmt.Sprintf("Карта в ltm сейчас *%s* и продлится *%dч %dм*\nСледующая карта *%s* и продлится *%dч %dм*",
						md2regex.ReplaceAllString(modes.Ltm.Current.Map, `\$1`),
						int(nextDiff.Hours()),
						int(nextDiff.Minutes())-int(nextDiff.Hours())*60,
						md2regex.ReplaceAllString(modes.Ltm.Next.Map, `\$1`),
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
