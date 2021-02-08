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
	"github.com/spf13/cobra"
)

const (
	defaultPubChannel string = "poll-result"
)

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
	}

	// Flags
	cmd.Flags().StringVar(&opts.redis.address, "redis-address", "", "the address where to connect to redis")
	cmd.Flags().StringVar(&opts.redis.pubChannel, "redis-pub-channel", defaultPubChannel, "redis channel where to subscribe from")

	cmd.MarkFlagRequired("redis-address")

	return cmd
}
