// Copyright © 2021 Elis Lulja
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	defaultPubChannel string = "poll-result"
)

var (
	log       zerolog.Logger
	locale, _ = time.LoadLocation("Europe/Rome")
)

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	}
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(locale)
	}
	log = zerolog.New(output).With().Timestamp().Logger().Level(zerolog.InfoLevel)
}

// NewRootCommand returns the root command
func NewRootCommand() *cobra.Command {
	opts := &options{
		redis: &redisOptions{},
	}

	cmd := &cobra.Command{
		Use:   "event-handler",
		Short: "handle events coming from the kube-scraper project",
		Long: `The event handler subscribes to events published by the Kube Scraper
project and reacts to them accirding to the channel and the event type.

In order to react to events, the redis address must be provided along with
the channel name.`,
		PreRun: func(_ *cobra.Command, _ []string) {
			if opts.debug {
				log.Level(zerolog.DebugLevel)
			}
		},
		Run: func(_ *cobra.Command, _ []string) {
			run(opts)
		},
	}

	// Flags
	cmd.Flags().StringVar(&opts.redis.address, "redis-address", "", "the address where to connect to redis")
	cmd.Flags().StringVar(&opts.redis.pubChannel, "redis-pub-channel", defaultPubChannel, "redis channel where to subscribe from")

	cmd.MarkFlagRequired("redis-address")

	return cmd
}

func run(opts *options) {
	// -- Init
	log.Info().Msg("starting...")
	ctx, canc := context.WithCancel(context.Background())
	exitChan := make(chan struct{})

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	go func() {
		defer close(exitChan)

		// -- Get redis client
		rdb, err := func() (*redis.Client, error) {
			_rdb := redis.NewClient(&redis.Options{Addr: opts.redis.address})
			rdCtx, rdCanc := context.WithTimeout(ctx, 15*time.Second)
			defer rdCanc()

			if _, err := _rdb.Ping(rdCtx).Result(); err != nil {
				log.Err(err).Msg("could not receive ping from redis, exiting...")
				return nil, err
			}

			return _rdb, nil
		}()
		if err != nil {
			signalChan <- os.Interrupt
			return
		}
		log.Info().Msg("connected to redis")
		defer rdb.Close()

		sub := rdb.Subscribe(ctx, opts.redis.pubChannel)
		defer sub.Close()

		iface, err := sub.Receive(ctx)
		if err != nil {
			log.Err(err).Str("channel", opts.redis.pubChannel).Msg("could not subscribe to channel")
			signalChan <- os.Interrupt
			return
		}
		l := log.With().Str("channel", opts.redis.pubChannel).Logger()

		switch iface.(type) {
		case *redis.Subscription:
			l.Info().Msg("subscribed to channel")
		case *redis.Message:
			go handleEvent(iface.(*redis.Message).Payload)
		case *redis.Pong:
			// pong received
		default:
			l.Error().Msg("error while getting subscription")
		}

		l.Info().Str("channel", opts.redis.pubChannel).Msg("listening for events...")
		ch := sub.Channel()
		select {
		case msg := <-ch:
			l.Info().Msg("received message")
			go handleEvent(msg.Payload)
		case <-ctx.Done():
			return
		}
	}()

	<-signalChan
	log.Info().Msg("exit requested")

	// -- Close all connections and shut down
	canc()
	<-exitChan

	log.Info().Msg("goodbye!")
}
