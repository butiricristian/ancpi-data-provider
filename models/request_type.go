package models

import "encoding/json"

type RequestType int64

const (
	ALTELE RequestType = iota
	INFORMARE
	INSCRIERE
	RECEPTIE
	UNDEFINED
)

func (rt RequestType) String() string {
	switch rt {
	case ALTELE:
		return "Altele"
	case INFORMARE:
		return "Informare"
	case INSCRIERE:
		return "Inscriere"
	case RECEPTIE:
		return "Receptie"
	}
	return "Unknown"
}

func (rt RequestType) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(rt.String())
	return value, err
}

func getRequestType(val string) RequestType {
	switch val {
	case "altele":
		return ALTELE
	case "informare":
		return INFORMARE
	case "inscriere":
		return INSCRIERE
	case "receptie":
		return RECEPTIE
	}
	return UNDEFINED
}
