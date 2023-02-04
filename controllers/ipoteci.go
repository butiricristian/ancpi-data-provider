package controllers

import (
	"fmt"
	"strconv"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
)

func filterIpoteciByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.IpoteciStateData {
	ipoteciData := map[time.Time][]*models.IpoteciStateData{}
	for _, val := range data.Data {
		if !dateStart.IsZero() && val.CurrentDate.Before(dateStart) {
			continue
		}
		if !dateEnd.IsZero() && val.CurrentDate.After(dateEnd) {
			continue
		}
		ipoteciData[val.CurrentDate] = val.IpoteciData
	}
	return ipoteciData
}

func filterIpoteciByJudet(result map[time.Time][]*models.IpoteciStateData, judet string) map[time.Time][]*models.IpoteciStateData {
	if judet == "" {
		judet = "TOTAL"
	}
	newResult := map[time.Time][]*models.IpoteciStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Name == judet {
				if _, ok := newResult[key]; !ok {
					newResult[key] = []*models.IpoteciStateData{}
				}
				newResult[key] = append(newResult[key], val)
			}
		}
	}
	return newResult
}

func filterIpoteciByActive(result map[time.Time][]*models.IpoteciStateData, active bool) map[time.Time]*models.IpoteciStateData {
	newResult := map[time.Time]*models.IpoteciStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Active == active {
				newResult[key] = val
			}
		}
	}
	return newResult
}

func HandleGetIpoteciData(judet string, activeString string, dateStartString string, dateEndString string) map[time.Time]*models.IpoteciStateData {
	fmt.Println("Getting Ipoteci Data")
	active, err := strconv.ParseBool(activeString)
	if err != nil {
		active = true
	}
	var dateStart time.Time
	var dateEnd time.Time
	if dateStartString != "" {
		dateStart = helpers.ConvertToTime(dateStartString)
	}
	if dateEndString != "" {
		dateEnd = helpers.ConvertToTime(dateEndString)
	}

	resultByInterval := filterIpoteciByInterval(dateStart, dateEnd)
	resultByJudet := filterIpoteciByJudet(resultByInterval, judet)
	result := filterIpoteciByActive(resultByJudet, active)
	return result
}
