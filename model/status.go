package model

import (
	"encoding/json"
	"errors"
)

type Status int

const (
	Created Status = iota
	Published
	Closed
	Canceled
)

func (s Status) String() string {
	switch s {
	case Created:
		return "Created"
	case Published:
		return "Published"
	case Closed:
		return "Closed"
	case Canceled:
		return "Canceled"
	default:
		return "Unknown"
	}
}

func ParseStatus(status string) (Status, error) {
	switch status {
	case "Created":
		return Created, nil
	case "Published":
		return Published, nil
	case "Closed":
		return Closed, nil
	case "Canceled":
		return Canceled, nil
	default:
		return Created, errors.New("invalid status")
	}
}
func (s *Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) (err error) {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}
	if *s, err = ParseStatus(status); err != nil {
		return err
	}
	return nil
}

func (s *Status) Scan(value interface{}) (err error) {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*s, err = ParseStatus(string(asBytes))
	return err
}
