package connection

import "github.com/streame-gg/go-discord-wrapper/types/events"

func (d *Client) OnEvent(
	eventName events.EventType,
	handler EventHandler,
) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.Events == nil {
		d.Events = make(map[events.EventType][]EventHandler)
	}

	d.Events[eventName] = append(d.Events[eventName], handler)
}

func (d *Client) OnGuildCreate(
	handler func(*Client, *events.GuildCreateEvent),
) {
	d.OnEvent(events.EventGuildCreate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.GuildCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to GuildCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnMessageCreate(
	handler func(*Client, *events.MessageCreateEvent),
) {
	d.OnEvent(events.EventMessageCreate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.MessageCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to MessageCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnInteractionCreate(
	handler func(*Client, *events.InteractionCreateEvent),
) {
	d.OnEvent(events.EventInteractionCreate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.InteractionCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to InteractionCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnReady(
	handler func(*Client, *events.ReadyEvent),
) {
	d.OnEvent(events.EventReady, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.ReadyEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to ReadyEvent: %T", event)
		}
	})
}

func (d *Client) OnGuildDelete(
	handler func(*Client, *events.GuildDeleteEvent),
) {
	d.OnEvent(events.EventGuildDelete, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.GuildDeleteEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to GuildDeleteEvent: %T", event)
		}
	})
}

func (d *Client) OnInviteCreate(
	handler func(*Client, *events.InviteCreateEvent),
) {
	d.OnEvent(events.EventInviteCreate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.InviteCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to InviteCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnInviteDelete(
	handler func(*Client, *events.InviteDeleteEvent),
) {
	d.OnEvent(events.EventInviteDelete, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.InviteDeleteEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to InviteDeleteEvent: %T", event)
		}
	})
}

func (d *Client) OnChannelCreate(
	handler func(*Client, *events.ChannelCreateEvent),
) {
	d.OnEvent(events.EventChannelCreate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.ChannelCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to ChannelCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnChannelDelete(
	handler func(*Client, *events.ChannelDeleteEvent),
) {
	d.OnEvent(events.EventChannelDelete, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.ChannelDeleteEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to ChannelDeleteEvent: %T", event)
		}
	})
}

func (d *Client) OnMessageDelete(
	handler func(*Client, *events.MessageDeleteEvent),
) {
	d.OnEvent(events.EventMessageDelete, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.MessageDeleteEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to MessageDeleteEvent: %T", event)
		}
	})
}

func (d *Client) OnMessageDeleteBulk(
	handler func(*Client, *events.MessageDeleteBulkEvent),
) {
	d.OnEvent(events.EventMessageDeleteBulk, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.MessageDeleteBulkEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to MessageDeleteBulkEvent: %T", event)
		}
	})
}

func (d *Client) OnMessageUpdate(
	handler func(*Client, *events.MessageUpdateEvent),
) {
	d.OnEvent(events.EventMessageUpdate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.MessageUpdateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to MessageUpdateEvent: %T", event)
		}
	})
}

func (d *Client) OnGuildAuditLogEntryCreate(
	handler func(*Client, *events.GuildAuditLogEntryCreateEvent),
) {
	d.OnEvent(events.EventGuildAuditLogEntryCreate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.GuildAuditLogEntryCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to GuildAuditLogEntryCreateEvent: %T", event)
		}
	})
}
