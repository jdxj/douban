package modules

type Opinion struct {
	ID     int
	Score  int
	Amount int
	One    float64
	Two    float64
	Three  float64
	Four   float64
	Five   float64
	Type   int
	Ref    int
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
	row = append(row, &o.Ref)
	return row
}
