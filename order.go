package angoslayer

import "fmt"

const (
	AnimeNameAsc  order = "anime_name_asc"
	AnimeNameDesc order = "anime_name_desc"
	AnimeYearAsc  order = "anime_year_asc"
	AnimeYearDesc order = "anime_year_desc"
	LatestFirst   order = "latest_first"
	EarlierFirst  order = "earliest_first"
	RatingDesc    order = "anime_rating_desc"
)

type order string

func (o order) valid() error {
	switch o {
	case AnimeNameAsc:
		fallthrough
	case AnimeNameDesc:
		fallthrough
	case AnimeYearAsc:
		fallthrough
	case AnimeYearDesc:
		fallthrough
	case LatestFirst:
		fallthrough
	case RatingDesc:
		fallthrough
	case EarlierFirst:
		return nil
	}
	return fmt.Errorf("Invalid order, Please use predefined orders by package")
}
