package transport

import (
	"github.com/devpablocristo/golang-monorepo/projects/tech-house/customer-manager-challenge/internal/customer/core/domain"
)

type GetKPIJson struct {
	AverageAge      float64 `json:"average_age"`
	AgeStdDeviation float64 `json:"age_std_deviation"`
}

func ToGetKPIJson(kpi *domain.KPI) *GetKPIJson {
	return &GetKPIJson{
		AverageAge:      kpi.AverageAge,
		AgeStdDeviation: kpi.AgeStdDeviation,
	}
}
