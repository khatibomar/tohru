package angoslayer

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type animeDetailsEndRes struct {
	Response AnimeDetails `json:"response"`
}

type moreInfoResult struct {
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
type commentFlagReasons struct {
	CommentFlagReasonID string `json:"comment_flag_reason_id"`
	FlagReason          string `json:"flag_reason"`
	FlagReasonOrder     string `json:"flag_reason_order"`
}
type contentRating struct {
	ContentType string `json:"content_type"`
	Level       string `json:"level"`
	VoteCount   string `json:"vote_count"`
}

type News struct {
	NewsID              string      `json:"news_id"`
	AnimeID             string      `json:"anime_id"`
	AnimeName           string      `json:"anime_name"`
	NewsTitle           string      `json:"news_title"`
	NewsDescription     string      `json:"news_description"`
	NewsDescriptionHTML interface{} `json:"news_description_html"`
	NewsImage           string      `json:"news_image"`
	NewsSourceLink      string      `json:"news_source_link"`
	NewsVideoLink       string      `json:"news_video_link"`
	NewsPublished       string      `json:"news_published"`
	NewsCreatedAt       string      `json:"news_created_at"`
	NewsUpdatedAt       interface{} `json:"news_updated_at"`
	NewsImageURL        string      `json:"news_image_url"`
	NewsUpdatedAtFormat string      `json:"news_updated_at_format"`
	NewsCreatedAtFormat string      `json:"news_created_at_format"`
}
type RelatedNews struct {
	Data []News `json:"data"`
}

type AnimeDetails struct {
	AnimeID              string               `json:"anime_id"`
	AnimeName            string               `json:"anime_name"`
	AnimeType            string               `json:"anime_type"`
	AnimeStatus          string               `json:"anime_status"`
	AnimeSeason          string               `json:"anime_season"`
	AnimeReleaseYear     string               `json:"anime_release_year"`
	AnimeAgeRating       string               `json:"anime_age_rating"`
	AnimeRating          string               `json:"anime_rating"`
	AnimeRatingUserCount string               `json:"anime_rating_user_count"`
	AnimeDescription     string               `json:"anime_description"`
	AnimeCoverImage      string               `json:"anime_cover_image"`
	AnimeTrailerURL      string               `json:"anime_trailer_url"`
	AnimeEnglishTitle    string               `json:"anime_english_title"`
	AnimeKeywords        string               `json:"anime_keywords"`
	AnimeUpdatedAt       string               `json:"anime_updated_at"`
	AnimeCreatedAt       string               `json:"anime_created_at"`
	AnimeGenreIds        string               `json:"anime_genre_ids"`
	AnimeGenres          string               `json:"anime_genres"`
	AnimeReleaseDay      string               `json:"anime_release_day"`
	AnimeCoverImageURL   string               `json:"anime_cover_image_url"`
	AnimeUpdatedAtFormat string               `json:"anime_updated_at_format"`
	AnimeCreatedAtFormat string               `json:"anime_created_at_format"`
	MoreInfoResult       moreInfoResult       `json:"more_info_result"`
	RelatedAnimes        interface{}          `json:"related_animes"`
	RelatedNews          interface{}          `json:"related_news"`
	CommentFlagReasons   []commentFlagReasons `json:"comment_flag_reasons"`
	ContentRating        []contentRating      `json:"content_rating"`
	Role                 string               `json:"role"`
}

type RelatedAnimes struct {
	Animes []Anime `json:"data"`
}

// UnmarshalJSON Implements a custom marsheler
// Because API have bad design and need to check edge cases
func (ad *AnimeDetails) UnmarshalJSON(b []byte) error {
	type details AnimeDetails
	var a details
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	switch a.RelatedAnimes.(type) {
	case []interface{}:
		a.RelatedAnimes = RelatedAnimes{}
	}
	switch a.RelatedNews.(type) {
	case []interface{}:
		a.RelatedNews = RelatedNews{}
	}
	*ad = AnimeDetails(a)
	return nil
}

func (s *AnimeService) GetAnimeDetails(animeID int) (AnimeDetails, error) {
	params := url.Values{}
	id := strconv.Itoa(animeID)
	params.Set("anime_id", id)
	params.Set("fetch_episodes", "No")
	params.Set("more_info", "Yes")

	res, err := s.getAnime(params, GetAnimeDetailsPath, http.MethodGet)
	if err != nil {
		return AnimeDetails{}, err
	}

	var details animeDetailsEndRes
	err = json.NewDecoder(res.Body).Decode(&details)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	return details.Response, err
}
