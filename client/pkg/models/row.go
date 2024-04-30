package models

type Row struct {
	Idx int
	Row   []string
}

func (r *Row) ToRow() []string {
	return r.Row
}

func (r *Row) Index() int {
	return r.Idx
}

func (r *Row) SetIndex(idx int){
	r.Idx = idx
}
