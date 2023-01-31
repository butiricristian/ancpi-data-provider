package models

import (
	"time"
)

type MonthlyData struct {
	CurrentDate time.Time
	VanzariData []*VanzariStateData
	CereriData  []*CereriStateData
	IpoteciData []*IpoteciStateData
}
