package connection

import (
	"context"
	"log/slog"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// webhookFetchTimeout bounds the REST fetch triggered by a WEBHOOKS_UPDATE.
const webhookFetchTimeout = 30 * time.Second

// hasWebhookHandlers reports whether any synthetic webhook handler is registered.
// The fetch + diff is skipped entirely when none are, so bots that don't use
// webhook events pay no REST cost.
func (d *Client) hasWebhookHandlers() bool {
	return len(d.discordEventEmitter.Handlers(events.EventWrapperWebhookCreate)) > 0 ||
		len(d.discordEventEmitter.Handlers(events.EventWrapperWebhookUpdate)) > 0 ||
		len(d.discordEventEmitter.Handlers(events.EventWrapperWebhookDelete)) > 0
}

// deriveWebhookSyntheticEventsAsync handles a WEBHOOKS_UPDATE off the gateway
// reader goroutine: WEBHOOKS_UPDATE carries only the channel, so the webhook list
// must be fetched over REST and diffed against the previous snapshot. It is a
// no-op when no webhook handler is registered.
func (d *Client) deriveWebhookSyntheticEventsAsync(ev *events.WebhooksUpdateEvent) {
	if !d.hasWebhookHandlers() {
		return
	}

	// Track the goroutine so Shutdown can drain it. Add before the shutdown check
	// to close the Add-after-Wait race (same pattern as enqueueOrDispatch).
	d.webhookWg.Add(1)
	select {
	case <-d.shutdownCh:
		d.webhookWg.Done()
		return
	default:
	}

	guildID, channelID := ev.GuildID, ev.ChannelID
	go func() {
		defer d.webhookWg.Done()
		defer func() {
			if r := recover(); r != nil {
				d.Logger.Error("panic in webhook synthetic derivation", slog.Any("recover", r))
			}
		}()
		d.fetchDiffEmitWebhooks(guildID, channelID)
	}()
}

// fetchDiffEmitWebhooks fetches the channel's current webhooks, diffs them
// against the cached snapshot, dispatches the resulting synthetic events, and
// updates the snapshot. On a cold snapshot it only seeds the baseline.
func (d *Client) fetchDiffEmitWebhooks(guildID, channelID discord.Snowflake) {
	ctx, cancel := context.WithTimeout(context.Background(), webhookFetchTimeout)
	defer cancel()

	// Abort the REST call promptly if the client shuts down mid-fetch.
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-d.shutdownCh:
			cancel()
		case <-done:
		}
	}()

	current, err := d.RestClient.ListChannelWebhooks(ctx, channelID)
	if err != nil {
		d.Logger.Debug("webhook synthetic: failed to list channel webhooks",
			slog.Any("channelId", channelID), slog.Any("err", err))
		return
	}
	for _, w := range current {
		w.Hydrate(d)
	}

	prev, hadSnapshot := d.loadWebhookSnapshot(channelID)
	d.storeWebhookSnapshot(channelID, current)

	// Cold snapshot: establish the baseline without firing (mirrors the
	// emoji/sticker cold-cache behaviour).
	if !hadSnapshot {
		return
	}

	for _, syn := range diffWebhooks(guildID, channelID, prev, current) {
		if d.enqueueOrDispatch(syn) {
			return
		}
	}
}

func (d *Client) loadWebhookSnapshot(channelID discord.Snowflake) ([]*discord.Webhook, bool) {
	v, ok := d.webhookSnapshots.Load(channelID.String())
	if !ok {
		return nil, false
	}
	webhooks, _ := v.([]*discord.Webhook)
	return webhooks, true
}

func (d *Client) storeWebhookSnapshot(channelID discord.Snowflake, webhooks []*discord.Webhook) {
	d.webhookSnapshots.Store(channelID.String(), webhooks)
}

// diffWebhooks compares the previous and current webhook sets for a channel and
// returns the synthetic Create/Update/Delete events to dispatch.
func diffWebhooks(guildID, channelID discord.Snowflake, old, current []*discord.Webhook) []events.Event {
	oldByID := make(map[discord.Snowflake]*discord.Webhook, len(old))
	for _, w := range old {
		oldByID[w.ID] = w
	}
	newByID := make(map[discord.Snowflake]*discord.Webhook, len(current))
	for _, w := range current {
		newByID[w.ID] = w
	}

	var result []events.Event

	for id, nw := range newByID {
		if ow, exists := oldByID[id]; !exists {
			result = append(result, &events.WebhookCreateEvent{
				GuildID:   guildID,
				ChannelID: channelID,
				Webhook:   nw,
			})
		} else if webhookChanged(ow, nw) {
			result = append(result, &events.WebhookUpdateEvent{
				GuildID:    guildID,
				ChannelID:  channelID,
				OldWebhook: ow,
				NewWebhook: nw,
			})
		}
	}
	for id, ow := range oldByID {
		if _, exists := newByID[id]; !exists {
			result = append(result, &events.WebhookDeleteEvent{
				GuildID:   guildID,
				ChannelID: channelID,
				Webhook:   ow,
			})
		}
	}

	return result
}

// webhookChanged reports whether the user-editable fields of a webhook differ.
func webhookChanged(old, new_ *discord.Webhook) bool {
	return !stringPtrEqual(old.Name, new_.Name) ||
		!stringPtrEqual(old.Avatar, new_.Avatar)
}
