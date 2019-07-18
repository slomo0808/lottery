package models

type ObjGiftPrize struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	PrizeNum     int    `json:"-"`
	LeftNum      int    `json:"-"`
	PrizeCodeA   int    `json:"-"`
	PrizeCodeB   int    `json:"-"`
	PrizeTime    int    `json:"-"`
	Img          string `json:"image"`
	Displayorder int    `json:"displayorder"`
	Gtype        int    `json:"gtype"`
	Gdata        string `json:"gdata"`
}
