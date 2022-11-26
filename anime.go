package tohru

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AnimeService service

const (
	PublishedAnimesPath = "anime/public/animes/get-published-animes"
	GetAnimeDetailsPath = "anime/public/anime/get-anime-details"
)

type animeEndRes struct {
	Response latestAnimeRespond `json:"response"`
}

type metaData struct {
	Limit   string `json:"_limit"`
	Offset  string `json:"_offset"`
	OrderBy string `json:"_order_by"`
}

type latestAnimeRespond struct {
	MetaData metaData `json:"meta_data"`
	Data     []Anime  `json:"data"`
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

func (s *AnimeService) GetLatestAnimes(offset, limit int) ([]Anime, error) {
	payload := JsonPayload{}
	var err error
	var payloadStr string

	err = payload.WithOffset(offset)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithLimit(limit)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithOrder(LatestFirst)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithListType(LatestUpdatedEpisodeNew)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithJustInfo("Yes")
	if err != nil {
		return []Anime{}, err
	}

	payloadStr, err = payload.ToJson()
	if err != nil {
		return []Anime{}, err
	}
	query := fmt.Sprintf(payloadStr, offset, limit)
	return s.getAnimeList(query)
}

func (s *AnimeService) SearchByName(offset, limit int, animeName string, orderBy order) ([]Anime, error) {
	payload := JsonPayload{}
	var err error
	var payloadStr string

	err = payload.WithOffset(offset)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithLimit(limit)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithOrder(orderBy)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithListType(Filter)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithJustInfo("Yes")
	if err != nil {
		return []Anime{}, err
	}
	payload.WithName(animeName)

	payloadStr, err = payload.ToJson()
	if err != nil {
		return []Anime{}, err
	}
	return s.getAnimeList(payloadStr)
}

func (s *AnimeService) OrderBy(offset, limit int, orderBy order) ([]Anime, error) {
	payload := JsonPayload{}
	var err error
	var payloadStr string

	err = payload.WithOffset(offset)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithLimit(limit)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithOrder(orderBy)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithListType(Filter)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithJustInfo("Yes")
	if err != nil {
		return []Anime{}, err
	}
	payload.WithName("")

	payloadStr, err = payload.ToJson()
	if err != nil {
		return []Anime{}, err
	}
	return s.getAnimeList(payloadStr)
}

func (s *AnimeService) GetAnimeListBySeason(offset, limit int, season season, orderBy order, releaseYear int) ([]Anime, error) {
	payload := JsonPayload{}
	var err error
	var payloadStr string

	err = payload.WithOffset(offset)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithLimit(limit)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithOrder(orderBy)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithListType(Filter)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithJustInfo("Yes")
	if err != nil {
		return []Anime{}, err
	}
	payload.WithName("")

	err = payload.WithSeason(season)
	if err != nil {
		return []Anime{}, err
	}
	err = payload.WithReleaseYear(releaseYear)
	if err != nil {
		return []Anime{}, err
	}

	payloadStr, err = payload.ToJson()
	if err != nil {
		return []Anime{}, err
	}
	return s.getAnimeList(payloadStr)
}

func (s *AnimeService) CustomAnimePayload(payload JsonPayload) ([]Anime, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return []Anime{}, err
	}
	return s.getAnimeList(string(data))
}

func (s *AnimeService) getAnime(params url.Values, path, method string) (*http.Response, error) {
	return s.getAnimeWithContext(context.Background(), params, path, method)
}

func (s *AnimeService) getAnimeWithContext(ctx context.Context, params url.Values, path, method string) (*http.Response, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = path

	u.RawQuery = params.Encode()

	var res *http.Response
	res, err := s.client.request(ctx, method, u.String(), nil)
	return res, err
}

func (s *AnimeService) getAnimeList(query string) ([]Anime, error) {
	params := url.Values{}
	params.Set("json", query)
	res, err := s.getAnime(params, PublishedAnimesPath, http.MethodGet)
	if err != nil {
		return []Anime{}, err
	}

	var animes animeEndRes
	err = json.NewDecoder(res.Body).Decode(&animes)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	return animes.Response.Data, err
}
