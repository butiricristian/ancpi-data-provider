package models

import (
	"fmt"

	"com.butiricristian/ancpi-data-provider/helpers"
)

type VanzariStateData struct {
	Name               string
	Agricol            int
	Neagricol          int
	Constructie        int
	FaraConstructie    int
	UnitatiIndividuale int
	Total              int
}

func (data *VanzariStateData) printData() string {
	return fmt.Sprintf("%v", *data)
}

func CreateVanzariData(row []string) VanzariStateData {
	return VanzariStateData{
		Name:               helpers.ReplaceSpecialCharacters(row[1]),
		Agricol:            helpers.ConvertToInt(row[2]),
		Neagricol:          helpers.ConvertToInt(row[3]),
		Constructie:        helpers.ConvertToInt(row[4]),
		FaraConstructie:    helpers.ConvertToInt(row[5]),
		UnitatiIndividuale: helpers.ConvertToInt(row[6]),
		Total:              helpers.ConvertToInt(row[7]),
	}
}
