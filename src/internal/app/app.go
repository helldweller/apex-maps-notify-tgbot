package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/sync/errgroup"

	"package/main/internal/apexapi"
)

var store modesStore

var ctx, cancel = context.WithCancel(context.Background())
var group, groupCtx = errgroup.WithContext(ctx)

const (
	defaultHTTPTimeout     = 30 * time.Second
	telegramPollingTimeout = 60 * time.Second
	telegramClientTimeout  = telegramPollingTimeout + 10*time.Second
)

func getExternalIP(httpClient *http.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://ifconfig.io/ip", nil)
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", fmt.Errorf("empty IP response")
	}

	return ip, nil
}

// Run func is similar to the main func.
func Run() {

	conf, err := loadConfig()
	if err != nil {
		fmt.Printf("Something went wrong while reading the configuration: %s", err)
		os.Exit(1)
	}
	configureLogger(conf)

	log.Info("Starting app")
	transport := http.DefaultTransport.(*http.Transport).Clone()
	httpClient := &http.Client{Timeout: defaultHTTPTimeout, Transport: transport}
	tgHTTPClient := &http.Client{Timeout: telegramClientTimeout, Transport: transport}

	if conf.HTTPProxy != "" {
		proxyURL, err := url.Parse(conf.HTTPProxy)
		if err != nil {
			log.Fatalf("Invalid HTTP proxy URL: %v", err)
		}

		transport.Proxy = http.ProxyURL(proxyURL)
	}

	apexapi.SetClient(httpClient)

	externalIP, err := getExternalIP(httpClient)
	if err != nil {
		log.Errorf("Failed to detect external IP: %v", err)
	} else {
		log.Infof("Application external IP: %s", externalIP)
	}

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
		// refresh fetches a fresh rotation without holding the store lock during
		// network I/O, then atomically swaps it in.
		refresh := func() {
			var fresh apexapi.Modes
			if err := fresh.Update(conf.ApexAPIKey); err != nil {
				log.Errorf("Failed to update map rotation: %v", err)
				return
			}
			store.set(fresh)
		}

		refresh() // get on start and then every UpdateInterval
		ticker := time.NewTicker(conf.UpdateInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				refresh()
			case <-groupCtx.Done():
				log.Error("Closing apexmaps modes goroutine")
				return groupCtx.Err()
			}
		}
	})

	bot, err := tgbotapi.NewBotAPIWithClient(conf.BotAPIKey, tgbotapi.APIEndpoint, tgHTTPClient)
	if err != nil {
		log.Panic(err)
	}
	tgbotapi.SetLogger(log)
	bot.Debug = conf.BotDebug
	log.Infof("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = int(telegramPollingTimeout / time.Second)
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

				command := update.Message.Command()
				switch command {
				case "help":
					msg.Text = "I understand /map /ranked /ltm"
				case "map":
					msg.Text = formatMode("", store.get().Pub, time.Now(), true)
				case "ranked":
					msg.Text = formatMode("в рейтинге ", store.get().Ranked, time.Now(), false)
				case "ltm":
					msg.Text = formatMode("в ltm ", store.get().Ltm, time.Now(), false)
				default:
					msg.Text = "I don't know that command, use /help"
				}

				// The map/ranked/ltm replies are formatted MarkdownV2 answers to
				// the triggering message; help/unknown are plain replies.
				if command == "map" || command == "ranked" || command == "ltm" {
					msg.ReplyToMessageID = update.Message.MessageID
					msg.ParseMode = "MarkdownV2"
					log.Infof("Recived new message from user %s, chat ID %d", update.Message.From, update.Message.Chat.ID)
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
