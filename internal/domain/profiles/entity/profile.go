package entity

type Profile struct {
	ID        uint64
	UserID    uint64
	FirstName string
	LastName  string
	Age       uint8
	Sex       Sex

	About    string
	Location string
}
