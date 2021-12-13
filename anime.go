package angoslayer

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
	query := fmt.Sprintf(`{"_offset":%d,"_limit":%d,"_order_by":"latest_first","list_type":"latest_updated_episode_new","just_info":"Yes"}`, offset, limit)
	return s.getAnimeList(query)
}

func (s *AnimeService) SearchByName(offset, limit int, animeName string, orderBy order) ([]Anime, error) {
	if err := orderBy.valid(); err != nil {
		return []Anime{}, err
	}
	query := fmt.Sprintf(`{"_offset":%d,"_limit":%d,"_order_by":"%s","list_type":"filter","anime_name":"%s","just_info":"Yes"}`, offset, limit, orderBy, animeName)
	return s.getAnimeList(query)
}

func (s *AnimeService) OrderBy(offset, limit int, orderBy order) ([]Anime, error) {
	if err := orderBy.valid(); err != nil {
		return []Anime{}, err
	}
	query := fmt.Sprintf(`{"_offset":%d,"_limit":%d,"_order_by":"%s","list_type":"filter","anime_name":"","just_info":"Yes"}`, offset, limit, orderBy)
	return s.getAnimeList(query)
}

func (s *AnimeService) GetAnimeListBySeason(offset, limit int, season season, orderBy order, releaseYear int) ([]Anime, error) {
	if err := orderBy.valid(); err != nil {
		return []Anime{}, err
	}
	if err := season.valid(); err != nil {
		return []Anime{}, err
	}
	query := fmt.Sprintf(`{"_offset":0,"_limit":30,"_order_by":"%s","list_type":"filter","anime_release_years":%d,"anime_season":"%s","just_info":"Yes"}`, orderBy, releaseYear, season)
	return s.getAnimeList(query)
}

func (s *AnimeService) CustomAnimePayload(payload string) ([]Anime, error) {
	return s.getAnimeList(payload)
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
