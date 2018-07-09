package models

import (
	"github.com/satori/go.uuid"
	. "github.com/Nastya-Kruglikova/cool_tasks/src/database"
	"net/url"
	sq "github.com/Masterminds/squirrel"
	"fmt"
	"errors"
)

const (
	addMuseumToTrip  = "INSERT INTO trips_museums (museum_id, trip_id) VALUES ($1, $2)"
	getMuseumsByTrip = "SELECT * FROM museums INNER JOIN trips_museums ON museums.id=trips_museums.museum_id AND trips_museums.trip_id=$1"
)

type Museum struct {
	ID         uuid.UUID
	Name       string
	Location   string
	Price      int
	OpenedAt   int
	ClosedAt   int
	MuseumType string
	Info       string
}

var AddMuseumToTrip = func(museum_id uuid.UUID, trip_id uuid.UUID) (error) {
	_, err := DB.Exec(addMuseumToTrip, museum_id, trip_id)
	return err
}

var GetMuseumsByTrip = func(trip_id uuid.UUID) ([]Museum, error) {
	rows, err := DB.Query(getMuseumsByTrip, trip_id)
	if err != nil {
		return []Museum{}, err
	}
	museums := make([]Museum, 0)

	for rows.Next() {
		var m Museum
		if err := rows.Scan(&m.ID, &m.Name, &m.Location, &m.Price, &m.OpenedAt, &m.ClosedAt, &m.MuseumType, &m.Info); err != nil {
			return []Museum{}, err
		}
		museums = append(museums, m)
	}
	return museums, nil
}

var GetMuseumsByRequest = func(params url.Values) ([]Museum, error) {
	museums := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("*").From("museums")
	var request string
	var err error
	var b sq.And
	count := 0
	validKeys := []string{"id", "name", "location", "price", "museum_type", "opened_at", "closed_at"}
	for key, value := range params {
		for _, keys := range validKeys {
			if key == keys {
				count++
			}
		}
		if count == 0 {
			return []Museum{}, errors.New("ERROR: Bad request")
		}
		switch key {
		case "name", "location", "museum_type":
			if len(value) > 1 {
				var or sq.Or
				for _, v := range value {
					or = append(or, sq.Eq{key: v})
				}
				b = append(b, or)
				request += ") AND "
			} else {
				b = append(b, sq.Eq{key: value[0]})
			}
		case "price", "opened_at", "closed_at":
			if len(value) == 2 {
				b = append(b, sq.And{sq.GtOrEq{key: value[1]}, sq.LtOrEq{key: value[0]}})
			} else {
				b = append(b, sq.Eq{key: value[0]})

			}
		default:
			b = append(b, sq.Eq{key: value[0]})
		}

		count = 0
	}

	request, _, err = museums.Where(b).ToSql()
	if err != nil {
		fmt.Println(err)
	}
	b = nil
	rows, err := DB.Query(request)
	if err != nil {
		return []Museum{}, err
	}

	museumsArray := make([]Museum, 0)

	for rows.Next() {
		var m Museum
		if err := rows.Scan(&m.ID, &m.Name, &m.Location, &m.Price, &m.OpenedAt, &m.ClosedAt, &m.MuseumType, &m.Info); err != nil {
			return []Museum{}, err
		}
		museumsArray = append(museumsArray, m)
	}
	return museumsArray, nil
}
