package model

import (
	"encoding/json"
	"errors"
)

type AuthorType int

const (
	Organization AuthorType = iota
	User
)

func (s AuthorType) String() string {
	switch s {
	case Organization:
		return "Organization"
	case User:
		return "User"
	default:
		return "Unknown"
	}
}

func ParseAuthorType(status string) (AuthorType, error) {
	switch status {
	case "Organization":
		return Organization, nil
	case "User":
		return User, nil
	default:
		return Organization, errors.New("invalid status")
	}
}
func (s *AuthorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *AuthorType) UnmarshalJSON(data []byte) (err error) {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}
	if *s, err = ParseAuthorType(status); err != nil {
		return err
	}
	return nil
}

func (s *AuthorType) Scan(value interface{}) (err error) {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*s, err = ParseAuthorType(string(asBytes))
	return err
}
