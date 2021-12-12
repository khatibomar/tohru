package angoslayer

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	LatestsAnimesPath = "anime/public/animes/get-published-animes"
)

type EndRes struct {
	Response Response `json:"response"`
}

func (r *EndRes) GetResult() []Anime {
	return r.Response.Data
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
	AnimeCoverImageURL string `json:"anime_cover_image_url"`
	AnimeTrailerURL    string `json:"anime_trailer_url"`
	AnimeReleaseDay    string `json:"anime_release_day"`
}

type Response struct {
	MetaData MetaData `json:"meta_data"`
	Data     []Anime  `json:"data"`
}

type AnimeService service

func (s *AnimeService) GetAnime(params url.Values) (*EndRes, error) {
	return s.GetAnimeListContext(context.Background(), params)
}

func (s *AnimeService) GetAnimeListContext(ctx context.Context, params url.Values) (*EndRes, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = LatestsAnimesPath

	u.RawQuery = params.Encode()

	var res EndRes
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	return &res, err
}

func (s *AnimeService) GetLatestAnimes(offset, limit int) ([]Anime, error) {
	params := url.Values{}
	query := fmt.Sprintf(`{"_offset":%d,"_limit":%d,"_order_by":"latest_first","list_type":"latest_updated_episode_new","just_info":"Yes"}`, offset, limit)
	params.Set("json", query)
	res, err := s.GetAnime(params)
	if err != nil {
		return []Anime{}, err
	}
	return res.GetResult(), nil
}
