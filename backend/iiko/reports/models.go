package reports

import "encoding/xml"

type Report struct {
	XMLName xml.Name `xml:"report"`
	Records []Record `xml:"r"`
}

// Record represents each <r> element in the XML.
type Record struct {
	OpenDateTime string  `xml:"OpenDate.Typed"`
	CloseTime    string  `xml:"CloseTime"`
	OpenTime     string  `xml:"OpenTime"`
	Department   string  `xml:"Department"`
	OrderNum     string  `xml:"OrderNum"`
	FullSum      float64 `xml:"fullSum"`
	CookingPlace string  `xml:"CookingPlace"`
	DishName     string  `xml:"DishName"`
}
