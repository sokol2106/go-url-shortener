package storage

type Shortdata struct {
	url   string
	short string
}

func NewShortdata(url string, short string) *Shortdata {
	return &Shortdata{url, short}
}

func (sd *Shortdata) Url() string {
	return sd.url
}

func (sd *Shortdata) Short() string {
	return sd.short
}
