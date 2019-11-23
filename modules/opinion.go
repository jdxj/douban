package modules

type Opinion struct {
	ID     int
	Score  float64
	Amount int64
	Type   int
	Ref    *int64 // 引用 Book.ID

	One, Two, Three float64
	Four, Five      float64
}

func (o *Opinion) ToScan() []interface{} {
	var row []interface{}
	row = append(row, &o.ID)
	row = append(row, &o.Score)
	row = append(row, &o.Amount)
	row = append(row, &o.One)
	row = append(row, &o.Two)
	row = append(row, &o.Three)
	row = append(row, &o.Four)
	row = append(row, &o.Five)
	row = append(row, &o.Type)
	row = append(row, o.Ref)
	return row
}

func (o *Opinion) ToInsert() []interface{} {
	var row []interface{}
	row = append(row, &o.Score)
	row = append(row, &o.Amount)
	row = append(row, &o.One)
	row = append(row, &o.Two)
	row = append(row, &o.Three)
	row = append(row, &o.Four)
	row = append(row, &o.Five)
	row = append(row, &o.Type)
	row = append(row, o.Ref)
	return row
}
