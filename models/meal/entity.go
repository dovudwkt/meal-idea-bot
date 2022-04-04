package meal

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"
)

const (
	sectionPhoto        = "photo:"
	sectionName         = "name:"
	sectionInstructions = "instructions:"
	sectionDescription  = "description:"

	sectionDelim = "==="
)

var (
	ErrNameEmpty  = errors.New("name not provided")
	ErrPhotoEmpty = errors.New("photo url empty or invalid")
)

type sqlMeal struct {
	ID           sql.NullInt64
	Name         sql.NullString
	PhotoURL     sql.NullString
	Description  sql.NullString
	Instructions sql.NullString
	CreatedAt    sql.NullInt64
}

func (m *sqlMeal) fromSQL() *Meal {
	meal := &Meal{
		ID:           m.ID.Int64,
		Name:         m.Name.String,
		PhotoURL:     m.PhotoURL.String,
		Description:  m.Description.String,
		Instructions: m.Instructions.String,
		CreatedAt:    m.CreatedAt.Int64,
	}

	return meal
}

type Meal struct {
	ID           int64
	Name         string
	PhotoURL     string
	Description  string
	Instructions string
	CreatedAt    int64
}

// ParseText parses the string into Meal struct. The string must have the following format:
//
//  photo: <photo url>
//  ===
//  name: <meal name>
//  ===
//  instructions: <cooling instructions>
//  ===
//  description: <description>
// Where === is a delimeter which has to be provided except after the last section.
func ParseText(text string) (*Meal, error) {
	meal := Meal{}

	sections := []string{sectionPhoto, sectionName, sectionInstructions, sectionDescription}

	paragraphs := strings.Split(text, sectionDelim)

	processedSections := make(map[string]struct{})

	for i, j := 0, 0; j < len(sections) && i < len(paragraphs); j++ {
		if _, processed := processedSections[sections[j]]; processed {
			i++
		}

		idx := strings.Index(paragraphs[i], sections[j])
		if idx == -1 {
			continue
		}

		value := paragraphs[i][idx+len(sections[j]):]
		value = strings.TrimSpace(value)

		switch sections[j] {
		case sectionPhoto:
			meal.PhotoURL = value
		case sectionName:
			meal.Name = value
		case sectionInstructions:
			meal.Instructions = value
		case sectionDescription:
			meal.Description = value
		}

		processedSections[sections[j]] = struct{}{}
		i++
	}

	if err := meal.Validate(); err != nil {
		return nil, err
	}

	return &meal, nil
}

func (m Meal) Validate() error {
	if m.Name == "" {
		return ErrNameEmpty
	}
	if m.PhotoURL == "" || !isURL(m.PhotoURL) {
		return ErrPhotoEmpty
	}

	return nil
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
