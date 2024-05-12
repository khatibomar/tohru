package tohru

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
	case AnimeNameAsc, AnimeNameDesc, AnimeYearAsc, AnimeYearDesc, LatestFirst, RatingDesc, EarlierFirst:
		return nil
	default:
		return fmt.Errorf("invalid order, Please use predefined orders by package")
	}
}
