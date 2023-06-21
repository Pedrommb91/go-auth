package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/rs/zerolog"
)

type queryBuilder[T any] struct {
	db       *sql.DB
	queryStr string
	mapper   QueryMapper[T]
}

type QueryMapper[T any] interface {
	Map(rows *sql.Rows) (T, error)
}

func With[T any](db *sql.DB) *queryBuilder[T] {
	b := &queryBuilder[T]{}
	b.db = db
	return b
}

func (q *queryBuilder[T]) Create(model T) (int64, error) {
	const op errors.Op = "database.Create"

	parser := NewModelParser[T](model)
	if parser.HasRelations() {
		tx, err := q.db.BeginTx(context.TODO(), nil)
		if err != nil {
			return 0, errors.Build(
				errors.WithOp(op),
				errors.WithError(err),
				errors.WithMessage("Failed to insert entry"),
			)
		}
		defer tx.Rollback()

		id, err := q.createWithRelations(tx, model)
		if err != nil {
			return 0, errors.Build(
				errors.WithOp(op),
				errors.WithError(err),
				errors.WithMessage("Failed to insert entry"),
			)
		}

		tx.Commit()
		return id, nil
	} else {
		return q.createWithoutRelations(model)
	}
}

func (q *queryBuilder[T]) Select(field string) *queryBuilder[T] {
	q.queryStr += " SELECT " + field
	return q
}

func (q *queryBuilder[T]) From(table string) *queryBuilder[T] {
	q.queryStr += " FROM " + table
	return q
}

func (q *queryBuilder[T]) WithMapper(mapper QueryMapper[T]) *queryBuilder[T] {
	q.mapper = mapper
	return q
}

func (q *queryBuilder[T]) Run() ([]T, error) {
	const op errors.Op = "database.Run"
	if !q.queryIsValid() {
		return nil, errors.Build(
			errors.WithOp(op),
			errors.WithError(fmt.Errorf("query is not valid")),
			errors.WithMessage("Query to database invalid"),
			errors.KindBadRequest(),
		)
	}

	rows, err := q.db.Query(q.queryStr)
	if err != nil {
		return nil, errors.Build(
			errors.WithOp(op),
			errors.WithMessage("Failed to get users from database"),
			errors.WithError(err),
			errors.WithSeverity(zerolog.ErrorLevel),
		)
	}

	data := make([]T, 0)
	for rows.Next() {
		element, err := q.mapper.Map(rows)
		if err != nil {
			return nil, errors.Build(
				errors.WithOp(op),
				errors.WithNestedErrorCopy(err),
			)
		}
		data = append(data, element)
	}
	return data, nil
}

func (q *queryBuilder[T]) queryIsValid() bool {
	return q.db != nil && q.queryStr != "" && q.mapper != nil
}

func (q *queryBuilder[T]) createWithoutRelations(model any) (int64, error) {
	const op errors.Op = "database.createWithoutRelations"

	parser := NewModelParser[T](model)

	sqlStatement := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s) RETURNING id;",
		parser.GetTableName(),
		parser.GetQueryColumns(),
		parser.GetQueryValues(),
	)

	var id int64
	err := q.db.QueryRowContext(context.TODO(), sqlStatement).Scan(&id)
	if err != nil {
		return 0, errors.Build(
			errors.WithOp(op),
			errors.WithMessage("Failed to insert entry"),
			errors.WithError(err),
			errors.WithSeverity(zerolog.ErrorLevel),
		)
	}

	return id, nil
}

func (q *queryBuilder[T]) createWithRelations(tx *sql.Tx, model any) (int64, error) {
	const op errors.Op = "database.createWithRelations"

	parser := NewModelParser[T](model)
	var references map[string]string = make(map[string]string)
	for _, v := range parser.GetAllRelationalStructs() {
		id, err := q.createWithRelations(tx, v)
		if err != nil {
			return 0, errors.Build(
				errors.WithOp(op),
				errors.WithError(err),
				errors.WithMessage("Failed to insert entry"),
			)
		}
		field := reflect.TypeOf(v).Name()
		references[parser.GetTagNameByTypeName(field)] = strconv.FormatInt(id, 10)
	}

	id, err := q.createWithParentRelations(tx, model, references)
	if err != nil {
		return 0, errors.Build(
			errors.WithOp(op),
			errors.WithError(err),
			errors.WithMessage("Failed to insert entry"),
		)
	}

	return id, nil
}

func (q *queryBuilder[T]) createWithParentRelations(tx *sql.Tx, model any, references map[string]string) (int64, error) {
	const op errors.Op = "database.createWithParentRelations"

	parser := NewModelParser[T](model)
	columns := parser.GetColumns()
	values := parser.GetValues()
	for i, v := range columns {
		if references[v] != "" {
			values[i] = references[v]
		}
	}

	sqlStatement := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s) RETURNING id;",
		parser.GetTableName(),
		strings.Join(columns[:], ", "),
		strings.Join(values[:], ", "),
	)

	var id int64
	err := tx.QueryRowContext(context.TODO(), sqlStatement).Scan(&id)
	if err != nil {
		return 0, errors.Build(
			errors.WithOp(op),
			errors.WithMessage("Failed to insert entry"),
			errors.WithError(err),
			errors.WithSeverity(zerolog.ErrorLevel),
		)
	}

	return id, nil
}
