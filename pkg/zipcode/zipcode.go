package zipcode

import (
	"regexp"

	"go-a-b-microservices/pkg/apperror"
)

type ZipCodeRequest struct {
	CEP string `json:"cep"`
}

func (z *ZipCodeRequest) Validate() error {
	if z.CEP == "" {
		return apperror.ErrZipCodeRequired
	}

	matched, _ := regexp.MatchString(`^\d{8}$`, z.CEP)
	if !matched {
		return apperror.ErrZipCodeInvalid
	}

	return nil
}

type Location struct {
	City string `json:"localidade"`
	CEP  string `json:"cep"`
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherData struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}
