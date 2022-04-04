package telegram

import (
	"errors"

	"meal-idea-bot/clients/telegram"
	"meal-idea-bot/events"
	meal_model "meal-idea-bot/models/meal"
)

type EventHandler struct {
	tg        *telegram.Client
	offset    int
	mealModel meal_model.Model
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, model meal_model.Model) *EventHandler {
	return &EventHandler{
		tg:        client,
		mealModel: model,
	}
}

func (p *EventHandler) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errors.New("error getting updates from Telegram: " + err.Error())
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

// Process handles the events and returns ErrUnknownEventType in case of invalid event type.
func (p *EventHandler) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

func (p *EventHandler) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errors.New("extract metadata from event: " + err.Error())
	}

	err = p.replyToCommand(event.Text, meta.ChatID, meta.Username)

	err = p.HandleClientErrors(err, meta.ChatID)
	if err != nil {
		return errors.New("reply to command: " + err.Error())
	}

	return nil
}

// meta extracts metadata from the event and returns ErrUnknownMetaType if event.Meta's type does not match Meta type.
func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMetaType
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
