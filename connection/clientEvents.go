package connection

import (
	"fmt"
	"reflect"

	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/util"
)

func (d *Client) onRawEvent(
	eventName events.EventType,
	handler EventHandler,
) {
	d.mu.Lock()
	if d.Events == nil {
		d.Events = make(map[events.EventType][]EventHandler)
	}
	d.Events[eventName] = append(d.Events[eventName], handler)
	d.mu.Unlock()

	if d.discordEventEmitter == nil {
		d.discordEventEmitter = util.NewEventEmitter[events.EventType, EventHandler]()
	}

	d.discordEventEmitter.On(eventName, handler)
}

func (d *Client) OnEvent(eventName events.EventType, handler interface{}) {
	wrapped, err := toEventHandler(handler, eventName)
	if err != nil {
		d.Logger.Warn().Msg(err.Error())
		return
	}

	d.onRawEvent(eventName, wrapped)
}

func toEventHandler(handler interface{}, eventName events.EventType) (EventHandler, error) {
	if typed, ok := handler.(EventHandler); ok {
		return typed, nil
	}

	handlerValue := reflect.ValueOf(handler)
	if !handlerValue.IsValid() || handlerValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("invalid handler for event %s: expected function", eventName)
	}

	handlerType := handlerValue.Type()
	if handlerType.NumIn() != 2 || handlerType.NumOut() != 0 {
		return nil, fmt.Errorf("invalid handler for event %s: expected func(*connection.Client, *events.X)", eventName)
	}

	clientType := reflect.TypeOf(&Client{})
	if !clientType.AssignableTo(handlerType.In(0)) {
		return nil, fmt.Errorf("invalid handler for event %s: first argument must be *connection.Client", eventName)
	}

	eventType := handlerType.In(1)

	return func(session *Client, event events.Event) {
		eventValue := reflect.ValueOf(event)
		if !eventValue.IsValid() || !eventValue.Type().AssignableTo(eventType) {
			session.Logger.Warn().Msgf(
				"Failed to cast event of type %T to expected type for event %s",
				event, eventName,
			)
			return
		}

		handlerValue.Call([]reflect.Value{reflect.ValueOf(session), eventValue})
	}, nil
}
