package book

type Book struct {
	ID     int
	Title  string
	Author string
	Press  string
}

func (b *Book) ToScan() []interface{} {
	var row []interface{}
	row = append(row, &b.ID)
	row = append(row, &b.Title)
	row = append(row, &b.Author)
	row = append(row, &b.Press)
	return row
}

func (b *Book) ToInsert() []interface{} {
	var row []interface{}
	row = append(row, &b.Title)
	row = append(row, &b.Author)
	row = append(row, &b.Press)
	return row
}
