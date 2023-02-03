package models

import (
	"fmt"

	"com.butiricristian/ancpi-data-provider/helpers"
)

type IpoteciStateData struct {
	VanzariStateData
	Active bool
}

func (data *IpoteciStateData) printData() string {
	return fmt.Sprintf("%v", *data)
}

func CreateIpoteciData(row []string) (IpoteciStateData, IpoteciStateData) {
	total := helpers.ConvertToInt(row[2]) + helpers.ConvertToInt(row[3]) + helpers.ConvertToInt(row[6]) + helpers.ConvertToInt(row[7])
	active := IpoteciStateData{
		VanzariStateData{
			Name:            helpers.ReplaceSpecialCharacters(row[1]),
			Agricol:         helpers.ConvertToInt(row[2]),
			Neagricol:       helpers.ConvertToInt(row[3]),
			Constructie:     helpers.ConvertToInt(row[6]),
			FaraConstructie: helpers.ConvertToInt(row[7]),
			Total:           total,
		},
		true,
	}
	total = helpers.ConvertToInt(row[4]) + helpers.ConvertToInt(row[5]) + helpers.ConvertToInt(row[8]) + helpers.ConvertToInt(row[9])
	inactive := IpoteciStateData{
		VanzariStateData{
			Name:            helpers.ReplaceSpecialCharacters(row[1]),
			Agricol:         helpers.ConvertToInt(row[4]),
			Neagricol:       helpers.ConvertToInt(row[5]),
			Constructie:     helpers.ConvertToInt(row[8]),
			FaraConstructie: helpers.ConvertToInt(row[9]),
			Total:           total,
		},
		false,
	}
	return active, inactive
}
