package data

type Error struct {
	Code string
	Desc string
}

func (e *Error) Error() string {
	return e.Desc
}
