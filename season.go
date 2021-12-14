package angoslayer

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
	case Fall:
		fallthrough
	case Winter:
		fallthrough
	case Summer:
		fallthrough
	case Spring:
		return nil
	}
	return fmt.Errorf("Invalid season, Please use predefined seasons by package")
}
