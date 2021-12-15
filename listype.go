package tohru

import "fmt"

const (
	CustomList              listType = "custom_list"
	AnimeList               listType = "anime_list"
	CurrentlyAiring         listType = "currently_airing"
	LatestUpdatedEpisode    listType = "latest_updated_episode"
	LatestUpdatedEpisodeNew listType = "latest_updated_episode_new"
	TopAnime                listType = "top_anime"
	TopCurrentlyAiring      listType = "top_currently_airing"
	TopTv                   listType = "top_tv"
	TopMovie                listType = "top_movie"
	Featured                listType = "featured"
	Filter                  listType = "filter"
	Favoirtes               listType = "watching"
	PlanToWatch             listType = "plan_to_watch"
	Watched                 listType = "watched"
	Dropped                 listType = "dropped"
	OnHold                  listType = "on_hold"
	WatchedHistory          listType = "watched_history"
	Schedule                listType = "schedule"
	LastAddedTv             listType = "last_added_tv"
	LastAddedMovie          listType = "last_added_movie"
	TopAnimeMal             listType = "top_anime_mal"
	CurrentlyAiringMal      listType = "top_currently_airing_mal"
	TopTvMal                listType = "top_tv_mal"
	AnimeCharacters         listType = "anime_characters"
	TopUpcoming             listType = "top_upcoming"
)

type listType string

func (l listType) valid() error {
	switch l {
	case CustomList:
		fallthrough
	case AnimeList:
		fallthrough
	case CurrentlyAiring:
		fallthrough
	case LatestUpdatedEpisode:
		fallthrough
	case LatestUpdatedEpisodeNew:
		fallthrough
	case TopAnime:
		fallthrough
	case TopCurrentlyAiring:
		fallthrough
	case TopTv:
		fallthrough
	case TopMovie:
		fallthrough
	case Featured:
		fallthrough
	case Filter:
		fallthrough
	case Favoirtes:
		fallthrough
	case PlanToWatch:
		fallthrough
	case Watched:
		fallthrough
	case Dropped:
		fallthrough
	case OnHold:
		fallthrough
	case WatchedHistory:
		fallthrough
	case Schedule:
		fallthrough
	case LastAddedTv:
		fallthrough
	case LastAddedMovie:
		fallthrough
	case TopAnimeMal:
		fallthrough
	case CurrentlyAiringMal:
		fallthrough
	case TopTvMal:
		fallthrough
	case AnimeCharacters:
		fallthrough
	case TopUpcoming:
		return nil
	}
	return fmt.Errorf("Invalid list type , Please use predefined list types by package")
}
