package entity

type Sex string

var (
	Female = Sex("female")
	Male   = Sex("male")

	ValidSex = map[Sex]bool{
		Female: true,
		Male:   true,
	}
)
