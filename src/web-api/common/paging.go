package common

// ======================== PhanTrang
type Paging struct {
	Page  int   `json:"Page" form:"page"`
	Limit int   `json:"Limit" form:"limit"`
	Total int64 `json:"Total" form:"-"`
}

func (p *Paging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit >= 100 {
		p.Limit = 10
	}
}
