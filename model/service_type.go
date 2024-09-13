package model

import (
	"encoding/json"
	"errors"
)

type ServiceType int

const (
	Construction ServiceType = iota
	Delivery
	Manufacture
)

func (s ServiceType) String() string {
	switch s {
	case Construction:
		return "Construction"
	case Delivery:
		return "Delivery"
	case Manufacture:
		return "Manufacture"
	default:
		return "Unknown"
	}
}

func ParseServiceType(s string) (ServiceType, error) {
	switch s {
	case "Construction":
		return Construction, nil
	case "Delivery":
		return Delivery, nil
	case "Manufacture":
		return Manufacture, nil
	default:
		return Construction, errors.New("invalid service type")
	}
}

func (s *ServiceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *ServiceType) UnmarshalJSON(data []byte) (err error) {
	var serviceType string
	if err := json.Unmarshal(data, &serviceType); err != nil {
		return err
	}
	if *s, err = ParseServiceType(serviceType); err != nil {
		return err
	}
	return nil
}

func (s *ServiceType) Scan(value interface{}) (err error) {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*s, err = ParseServiceType(string(asBytes))
	return err
}
