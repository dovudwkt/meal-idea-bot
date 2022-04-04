package meal

import (
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Model interface {
	Create(*Meal) (id int64, err error)
	GetByID(int64) (*Meal, error)
	Search(SearchOptions) ([]Meal, error)
	GetAll(limit int) ([]Meal, error)
	GetRandom() (*Meal, error)
	Delete(id int64) error
}

type SearchOptions struct {
	ID            int64
	Name          string
	CreatedAfter  int64
	CreatedBefore int64
}

type model struct {
	db *sql.DB
}

func New(driverName, dsn string) (Model, error) {
	db, err := sql.Open(driverName, dsn)
	return &model{db: db}, err
}

var sqlSelect = sq.Select("id", "name", "photo_url", "instructions", "description", "created_at").From("meal")

func (m *model) Create(meal *Meal) (int64, error) {
	q := "INSERT INTO meal(name, photo_url, instructions, description, created_at) VALUES (?,?,?,?,?)"

	sqlRes, err := m.db.Exec(q, meal.Name, meal.PhotoURL, meal.Instructions, meal.Description, time.Now().Unix())
	if err != nil {
		return 0, errors.New("insert sql: " + err.Error())
	}

	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, errors.New("last insert id: " + err.Error())
	}

	return id, nil
}

func (m *model) GetByID(id int64) (*Meal, error) {
	builder := sqlSelect.Where(sq.Eq{"id": id})

	q, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("to sql: " + err.Error())
	}

	var row sqlMeal

	err = m.db.QueryRow(q, args...).Scan(&row.ID, &row.Name, &row.PhotoURL, &row.Instructions, &row.Description, &row.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, errors.New("scan row: " + err.Error())
	}

	return row.fromSQL(), nil
}

func (m *model) Search(opts SearchOptions) ([]Meal, error) {
	sql := sqlSelect

	if opts.ID != 0 {
		sql = sql.Where(sq.Eq{"id": opts.ID})
	}
	if opts.Name != "" {
		sql = sql.Where(sq.Eq{"name": opts.Name})
	}
	if opts.CreatedBefore != 0 {
		sql = sql.Where(sq.Lt{"created_at": opts.CreatedBefore})
	}
	if opts.CreatedAfter != 0 {
		sql = sql.Where(sq.Lt{"created_at": opts.CreatedAfter})
	}

	q, args, err := sql.ToSql()
	if err != nil {
		return nil, errors.New("to sql: " + err.Error())
	}

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, errors.New("execute sql: " + err.Error())
	}
	defer rows.Close()

	var meals []Meal

	for rows.Next() {
		var row sqlMeal

		err = rows.Scan(&row.ID, &row.Name, &row.PhotoURL, &row.Instructions, &row.Description, &row.CreatedAt)
		if err != nil {
			return nil, errors.New("scan row: " + err.Error())
		}

		meals = append(meals, *row.fromSQL())
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("rows error: " + err.Error())
	}

	return meals, nil
}

func (m *model) GetAll(limit int) ([]Meal, error) {
	q, args, err := sqlSelect.Limit(uint64(limit)).ToSql()
	if err != nil {
		return nil, errors.New("to sql: " + err.Error())
	}

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, errors.New("execute sql: " + err.Error())
	}
	defer rows.Close()

	var meals []Meal

	for rows.Next() {
		var row sqlMeal

		err = rows.Scan(&row.ID, &row.Name, &row.PhotoURL, &row.Instructions, &row.Description, &row.CreatedAt)
		if err != nil {
			return nil, errors.New("scan row: " + err.Error())
		}

		meals = append(meals, *row.fromSQL())
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("rows error: " + err.Error())
	}

	return meals, nil
}

func (m *model) GetRandom() (*Meal, error) {
	builder := sqlSelect.OrderBy("RAND()").Limit(1)

	q, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("to sql: " + err.Error())
	}

	var row sqlMeal

	err = m.db.QueryRow(q, args...).Scan(&row.ID, &row.Name, &row.PhotoURL, &row.Instructions, &row.Description, &row.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, errors.New("scan row: " + err.Error())
	}

	return row.fromSQL(), nil
}

func (m *model) Delete(id int64) error {

	return nil
}
