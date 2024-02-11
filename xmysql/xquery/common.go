package xquery

const (
	defaultLimit = 50
)

type Page struct {
	Start int    `json:"start"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

func GetOrNewPage(page *Page) *Page {
	if page == nil {
		page = &Page{Limit: defaultLimit}
	}
	return page
}
