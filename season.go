package tohru

import "fmt"

const (
	Fall   season = "Fall"
	Summer season = "Summer"
	Winter season = "Winter"
	Spring season = "Spring"
)

type season string

func (s season) valid() error {
	switch s {
	case Fall, Winter, Summer, Spring:
		return nil
	default:
		return fmt.Errorf("invalid season, Please use predefined seasons by package")
	}
}
