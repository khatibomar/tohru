package angoslayer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type order string
type list string
type season string

func (o order) valid() error {
	switch o {
	case AnimeNameAsc:
	case AnimeNameDesc:
	case AnimeYearAsc:
	case AnimeYearDesc:
	case LatestFirst:
	case RatingDesc:
	case EarlierFirst:
		return nil
	}
	return fmt.Errorf("Invalid order, Please use predefined orders by package")
}

func (s season) valid() error {
	switch s {
	case Fall:
	case Winter:
	case Summer:
	case Spring:
		return nil
	}
	return fmt.Errorf("Invalid season, Please use predefined seasons by package")
}

func (l list) valid() error {
	switch l {
	case CustomList:
	case AnimeList:
	case CurrentlyAiring:
	case LatestUpdatedEpisode:
	case LatestUpdatedEpisodeNew:
	case TopAnime:
	case TopCurrentlyAiring:
	case TopTv:
	case TopMovie:
	case Featured:
	case Filter:
	case Favoirtes:
	case PlanToWatch:
	case Watched:
	case Dropped:
	case OnHold:
	case WatchedHistory:
	case Schedule:
	case LastAddedTv:
	case LastAddedMovie:
	case TopAnimeMal:
	case CurrentlyAiringMal:
	case TopTvMal:
	case AnimeCharacters:
	case TopUpcoming:
		return nil
	}
	return fmt.Errorf("Invalid list type , Please use predefined list types by package")
}

const (
	PublishedAnimesPath = "anime/public/animes/get-published-animes"
	GetAnimeDetailsPath = "anime/public/anime/get-anime-details"

	AnimeNameAsc  order = "anime_name_asc"
	AnimeNameDesc order = "anime_name_desc"
	AnimeYearAsc  order = "anime_year_asc"
	AnimeYearDesc order = "anime_year_desc"
	LatestFirst   order = "latest_first"
	EarlierFirst  order = "earliest_first"
	RatingDesc    order = "anime_rating_desc"

	CustomList              list = "custom_list"
	AnimeList               list = "anime_list"
	CurrentlyAiring         list = "currently_airing"
	LatestUpdatedEpisode    list = "latest_updated_episode"
	LatestUpdatedEpisodeNew list = "latest_updated_episode_new"
	TopAnime                list = "top_anime"
	TopCurrentlyAiring      list = "top_currently_airing"
	TopTv                   list = "top_tv"
	TopMovie                list = "top_movie"
	Featured                list = "featured"
	Filter                  list = "filter"
	Favoirtes               list = "watching"
	PlanToWatch             list = "plan_to_watch"
	Watched                 list = "watched"
	Dropped                 list = "dropped"
	OnHold                  list = "on_hold"
	WatchedHistory          list = "watched_history"
	Schedule                list = "schedule"
	LastAddedTv             list = "last_added_tv"
	LastAddedMovie          list = "last_added_movie"
	TopAnimeMal             list = "top_anime_mal"
	CurrentlyAiringMal      list = "top_currently_airing_mal"
	TopTvMal                list = "top_tv_mal"
	AnimeCharacters         list = "anime_characters"
	TopUpcoming             list = "top_upcoming"

	Fall   season = "Fall"
	Summer season = "Summer"
	Winter season = "Winter"
	Spring season = "Spring"
)

type AnimeEndRes struct {
	Response LatestAnimeRespond `json:"response"`
}

func (r *AnimeEndRes) GetResult() []Anime {
	return r.Response.Data
}

type AnimeDetailsEndRes struct {
	Response AnimeDetails `json:"response"`
}

func (r *AnimeDetailsEndRes) GetResult() AnimeDetails {
	return r.Response
}

type MoreInfoResult struct {
	Score          string      `json:"score"`
	ScoredBy       string      `json:"scored_by"`
	TrailerURL     string      `json:"trailer_url"`
	Source         string      `json:"source"`
	Episodes       interface{} `json:"episodes"`
	Duration       string      `json:"duration"`
	AiredFrom      string      `json:"aired_from"`
	AiredTo        interface{} `json:"aired_to"`
	AnimeStudioIds string      `json:"anime_studio_ids"`
	AnimeStudios   string      `json:"anime_studios"`
}
type CommentFlagReasons struct {
	CommentFlagReasonID string `json:"comment_flag_reason_id"`
	FlagReason          string `json:"flag_reason"`
	FlagReasonOrder     string `json:"flag_reason_order"`
}
type ContentRating struct {
	ContentType string `json:"content_type"`
	Level       string `json:"level"`
	VoteCount   string `json:"vote_count"`
}
type AnimeDetails struct {
	AnimeID                string               `json:"anime_id"`
	AnimeName              string               `json:"anime_name"`
	AnimeType              string               `json:"anime_type"`
	AnimeStatus            string               `json:"anime_status"`
	JustInfo               string               `json:"just_info"`
	AnimeFeatured          string               `json:"anime_featured"`
	AnimeSeason            string               `json:"anime_season"`
	AnimeReleaseYear       string               `json:"anime_release_year"`
	AnimeAgeRating         string               `json:"anime_age_rating"`
	AnimeRating            string               `json:"anime_rating"`
	AnimeRatingUserCount   string               `json:"anime_rating_user_count"`
	AnimeDescription       string               `json:"anime_description"`
	AnimeCoverImage        string               `json:"anime_cover_image"`
	AnimeCoverImageFull    string               `json:"anime_cover_image_full"`
	AnimeBannerImage       interface{}          `json:"anime_banner_image"`
	AnimeTrailerURL        string               `json:"anime_trailer_url"`
	AnimeEnglishTitle      string               `json:"anime_english_title"`
	AnimeKeywords          string               `json:"anime_keywords"`
	AllowComment           string               `json:"allow_comment"`
	AnimeUpdatedAt         string               `json:"anime_updated_at"`
	AnimeCreatedAt         string               `json:"anime_created_at"`
	AnimeGenreIds          string               `json:"anime_genre_ids"`
	AnimeGenres            string               `json:"anime_genres"`
	AnimeReleaseDay        string               `json:"anime_release_day"`
	AnimeCoverImageURL     string               `json:"anime_cover_image_url"`
	AnimeCoverImageFullURL string               `json:"anime_cover_image_full_url"`
	AnimeBannerImageURL    interface{}          `json:"anime_banner_image_url"`
	AnimeUpdatedAtFormat   string               `json:"anime_updated_at_format"`
	AnimeCreatedAtFormat   string               `json:"anime_created_at_format"`
	MoreInfoResult         MoreInfoResult       `json:"more_info_result"`
	RelatedAnimes          []interface{}        `json:"related_animes"`
	RelatedNews            []interface{}        `json:"related_news"`
	CommentFlagReasons     []CommentFlagReasons `json:"comment_flag_reasons"`
	ContentRating          []ContentRating      `json:"content_rating"`
	Role                   string               `json:"role"`
	TopAnimeContributors   []interface{}        `json:"top_anime_contributors"`
}

type MetaData struct {
	Limit   string `json:"_limit"`
	Offset  string `json:"_offset"`
	OrderBy string `json:"_order_by"`
}

type Anime struct {
	AnimeID            string `json:"anime_id"`
	AnimeName          string `json:"anime_name"`
	AnimeType          string `json:"anime_type"`
	AnimeStatus        string `json:"anime_status"`
	JustInfo           string `json:"just_info"`
	AnimeSeason        string `json:"anime_season"`
	AnimeReleaseYear   string `json:"anime_release_year"`
	AnimeRating        string `json:"anime_rating"`
	LatestEpisodeID    string `json:"latest_episode_id"`
	LatestEpisodeName  string `json:"latest_episode_name"`
	AnimeGenres        string `json:"anime_genres"`
	AnimeCoverImageURL string `json:"anime_cover_image_url"`
	AnimeTrailerURL    string `json:"anime_trailer_url"`
	AnimeReleaseDay    string `json:"anime_release_day"`
}

type LatestAnimeRespond struct {
	MetaData MetaData `json:"meta_data"`
	Data     []Anime  `json:"data"`
}

type AnimeService service

func (s *AnimeService) GetAnime(params url.Values, path, method string) (*http.Response, error) {
	return s.GetAnimeWithContext(context.Background(), params, path, method)
}

func (s *AnimeService) GetAnimeWithContext(ctx context.Context, params url.Values, path, method string) (*http.Response, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = path

	u.RawQuery = params.Encode()

	var res *http.Response
	res, err := s.client.Request(ctx, method, u.String(), nil)
	return res, err
}

func (s *AnimeService) GetAnimeList(offset, limit int, query string) ([]Anime, error) {
	params := url.Values{}
	params.Set("json", query)
	res, err := s.GetAnime(params, PublishedAnimesPath, http.MethodGet)
	if err != nil {
		return []Anime{}, err
	}

	var animes AnimeEndRes
	err = json.NewDecoder(res.Body).Decode(&animes)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	return animes.GetResult(), err
}

func (s *AnimeService) GetLatestAnimes(offset, limit int) ([]Anime, error) {
	query := fmt.Sprintf(`{"_offset":%d,"_limit":%d,"_order_by":"latest_first","list_type":"latest_updated_episode_new","just_info":"Yes"}`, offset, limit)
	return s.GetAnimeList(offset, limit, query)
}

func (s *AnimeService) SearchByName(offset, limit int, animeName string, orderBy order) ([]Anime, error) {
	if err := orderBy.valid(); err != nil {
		return []Anime{}, err
	}
	query := fmt.Sprintf(`{"_offset":%d,"_limit":%d,"_order_by":"%s","list_type":"filter","anime_name":"%s","just_info":"Yes"}`, offset, limit, orderBy, animeName)
	return s.GetAnimeList(offset, limit, query)
}

func (s *AnimeService) OrderBy(offset, limit int, orderBy order) ([]Anime, error) {
	if err := orderBy.valid(); err != nil {
		return []Anime{}, err
	}
	query := fmt.Sprintf(`{"_offset":%d,"_limit":%d,"_order_by":"%s","list_type":"filter","anime_name":"","just_info":"Yes"}`, offset, limit, orderBy)
	return s.GetAnimeList(offset, limit, query)
}

func (s *AnimeService) GetAnimeListBySeason(offset, limit int, season season, orderBy order, releaseYear int) ([]Anime, error) {
	if err := orderBy.valid(); err != nil {
		return []Anime{}, err
	}
	if err := season.valid(); err != nil {
		return []Anime{}, err
	}
	query := fmt.Sprintf(`https://anslayer.com/anime/public/animes/get-published-animes?json={"_offset":0,"_limit":30,"_order_by":"%s","list_type":"filter","anime_release_years":%d,"anime_season":"%s","just_info":"Yes"}`, orderBy, releaseYear, season)
	return s.GetAnimeList(offset, limit, query)
}

func (s *AnimeService) GetAnimeDetails(animeID int) (AnimeDetails, error) {
	params := url.Values{}
	id := strconv.Itoa(animeID)
	params.Set("anime_id", id)
	params.Set("fetch_episodes", "No")
	params.Set("more_info", "Yes")

	res, err := s.GetAnime(params, GetAnimeDetailsPath, http.MethodGet)
	if err != nil {
		return AnimeDetails{}, err
	}

	var details AnimeDetailsEndRes
	err = json.NewDecoder(res.Body).Decode(&details)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	return details.GetResult(), err
}
