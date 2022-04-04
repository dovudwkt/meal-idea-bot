package telegram

import (
	"errors"
	"log"
	"strings"

	meal_model "meal-idea-bot/models/meal"
)

const (
	AddMealCmd   = "/add"
	AddSampleCmd = "/add_sample"
	MealCmd      = "/meal"
	HelpCmd      = "/help"
	StartCmd     = "/start"
)

func (p *EventHandler) replyToCommand(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	switch text {
	case MealCmd:
		return p.sendRandomMeal(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case AddSampleCmd:
		return p.sendAddSample(chatID)
		// default:
		// 	return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
	if isAddMealCmd(text) {
		return p.AddMeal(chatID, text)
	}

	return p.tg.SendMessage(chatID, msgUnknownCommand)
}

// HandleClientErrors checks if the error chain contains any of the client errors.
// If so it will reply to the client with the error message, otherwise,
// same error will be returned back to check for other errors.
func (p *EventHandler) HandleClientErrors(err error, chatID int) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), meal_model.ErrNameEmpty.Error()) {
		return p.tg.SendMessage(chatID, meal_model.ErrNameEmpty.Error())
	}
	if strings.Contains(err.Error(), meal_model.ErrPhotoEmpty.Error()) {
		return p.tg.SendMessage(chatID, meal_model.ErrPhotoEmpty.Error())
	}

	return err
}

func (p *EventHandler) AddMeal(chatID int, text string) error {
	meal, err := meal_model.ParseText(text)
	if err != nil {
		return errors.New("error parsing text to meal struct: " + err.Error())
	}

	if _, err = p.mealModel.Create(meal); err != nil {
		return err
	}

	if err = p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	log.Printf("meal '%s' saved", meal.Name)

	return nil
}

func (p *EventHandler) sendRandomMeal(chatID int) (err error) {
	// meal := data.PickRandomMeal()
	meal, err := p.mealModel.GetRandom()
	if err != nil {
		return errors.New("get random meal: " + err.Error())
	}

	if err := p.tg.SendPhoto(chatID, meal.PhotoURL, meal.Name); err != nil {
		return errors.New("can't send photo: " + err.Error())
	}

	if meal.Description != "" {
		if err := p.tg.SendMessage(chatID, meal.Description); err != nil {
			return errors.New("can't send description: " + err.Error())
		}
	}
	if meal.Instructions != "" {
		if err := p.tg.SendMessage(chatID, meal.Instructions); err != nil {
			return errors.New("can't send instructions: " + err.Error())
		}
	}

	return nil
}

func (p *EventHandler) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *EventHandler) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *EventHandler) sendAddSample(chatID int) error {
	return p.tg.SendMessage(chatID, msgAddSample)
}

func isAddMealCmd(text string) bool {
	return strings.Contains(text, AddMealCmd)
}
