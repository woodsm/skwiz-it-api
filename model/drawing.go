package model

type Drawing struct {
	Id     int64   `json:"id"`
	Url    *string `json:"url"`
	Top    Section `json:"top"`
	Middle Section `json:"middle"`
	Bottom Section `json:"bottom"`
}
