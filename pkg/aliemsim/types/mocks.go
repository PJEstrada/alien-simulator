package types

import (
	"errors"
	"io"
)

var ErrorNewMapMock = errors.New("ErrorNewMapMock")

func NewMapFromReaderMock(reader io.Reader) (*Map[City, Direction], error) {
	return &Map[City, Direction]{
		Cities: CityStore{
			"test1": &City{Name: "test1"},
			"test2": &City{Name: "test2"},
			"test3": &City{Name: "test3"},
		},
	}, nil
}

func NewMapFromReaderMockErr(reader io.Reader) (*Map[City, Direction], error) {
	return &Map[City, Direction]{}, ErrorNewMapMock
}
