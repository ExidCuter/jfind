package carQuery

import (
	"reflect"
)

const URL string = "https://www.avto.net/Ads/results.asp"

type CarQuery struct {
	EQ1              string
	EQ2              string
	EQ3              string
	EQ4              string
	EQ5              string
	EQ6              string
	EQ7              string
	EQ8              string
	EQ9              string
	EToznaka         string
	KAT              string
	PIA              string
	PIAzero          string
	PSLO             string
	Airbag           string
	Akcija           string
	Arhiv            string
	Barva            string
	Barvaint         string
	Bencin           string
	Broker           string
	CcmMax           string
	CcmMin           string
	CenaMax          string
	CenaMin          string
	Col              string
	Dolzina          string
	DolzinaMAX       string
	DolzinaMIN       string
	Kategorija       string
	KmMax            string
	KmMin            string
	KwMax            string
	KwMin            string
	LetnikMax        string
	LetnikMin        string
	Lezisc           string
	Lokacija         string
	MocMax           string
	MocMin           string
	Model            string
	Model2           string
	Model3           string
	ModelID          string
	Motortakt        string
	Motorvalji       string
	NosilnostMAX     string
	NosilnostMIN     string
	Oblika           string
	Paketgarancije   string
	Premer           string
	Presek           string
	Presort          string
	Prikazkategorije string
	Sirina           string
	Starost2         string
	Stran            string
	SubSORT          string
	SubTIPSORT       string
	Tip              string
	Tip2             string
	Tip3             string
	Tipsort          string
	Vijakov          string
	Vozilo           string
	Zaloga           string
	Znamka           string
	Znamka2          string
	Znamka3          string
}

func (carq CarQuery) GetURI() string {
	uri := ""

	v := reflect.ValueOf(carq)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() != "" {
			uri += "&" + v.Type().Field(i).Name + "=" + v.Field(i).String()
		}
	}

	uri = uri[1:]

	return URL + "?" + uri
}
