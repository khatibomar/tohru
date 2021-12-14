package angoslayer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	GetEpisodePath      = "/anime/public/episodes/get-episodes-new"
	EpisodeDownloadPath = "/la/public/api/fw"
)

type EpisodeService service

type episodeEndRes struct {
	Response episodesResponse `json:"response"`
}

type episodeUrls struct {
	EpisodeURLID      string `json:"episode_url_id"`
	EpisodeServerID   string `json:"episode_server_id"`
	EpisodeServerName string `json:"episode_server_name"`
	EpisodeURL        string `json:"episode_url"`
}

type nextEpisode struct {
	EpisodeID   string `json:"episode_id"`
	EpisodeName string `json:"episode_name"`
}

type Episode struct {
	EpisodeID                 string        `json:"episode_id"`
	EpisodeName               string        `json:"episode_name"`
	EpisodeNumber             string        `json:"episode_number"`
	AllowComment              string        `json:"allow_comment"`
	SkipFrom                  string        `json:"skip_from"`
	SkipTo                    string        `json:"skip_to"`
	EpisodeRating             string        `json:"episode_rating"`
	EpisodeRatingUserCount    string        `json:"episode_rating_user_count"`
	EpisodeWatchedHistory     interface{}   `json:"episode_watched_history"`
	EpisodeAlreadyRatedByUser interface{}   `json:"episode_already_rated_by_user"`
	EpisodeRatingByUser       interface{}   `json:"episode_rating_by_user"`
	EpisodeUrls               []episodeUrls `json:"episode_urls"`
	NextEpisode               []nextEpisode `json:"next_episode"`
	PreviousEpisode           []interface{} `json:"previous_episode"`
}

type episodesResponse struct {
	Episodes []Episode `json:"data"`
	Count    int       `json:"count"`
}

type DownloadLinks []string

func (s *EpisodeService) getEpisode(params url.Values, path, method, payload string) (*http.Response, error) {
	return s.getEpisodeWithContext(context.Background(), params, path, method, payload)
}

func (s *EpisodeService) getEpisodeWithContext(ctx context.Context, params url.Values, path, method, payload string) (*http.Response, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = path

	u.RawQuery = params.Encode()

	payloadReader := strings.NewReader(payload)

	var res *http.Response
	res, err := s.client.request(ctx, method, u.String(), payloadReader)
	return res, err
}

func (s *EpisodeService) GetEpisodesList(animeID int) ([]Episode, error) {
	params := url.Values{}
	payload := fmt.Sprintf(`json={"more_info":"No","anime_id":%d}`, animeID)
	res, err := s.getEpisode(params, GetEpisodePath, http.MethodPost, payload)
	if err != nil {
		return []Episode{}, err
	}

	var episodes episodeEndRes
	err = json.NewDecoder(res.Body).Decode(&episodes)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	return episodes.Response.Episodes, err
}

func (s *EpisodeService) GetEpisodeDetails(animeID, episodeID int) (Episode, error) {
	params := url.Values{}
	payload := fmt.Sprintf(`json={"episode_id":%d,"anime_id":%d}`, episodeID, animeID)
	res, err := s.getEpisode(params, GetEpisodePath, http.MethodPost, payload)
	if err != nil {
		return Episode{}, err
	}

	var episodes episodeEndRes
	err = json.NewDecoder(res.Body).Decode(&episodes)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	return episodes.Response.Episodes[0], err
}

func (s *EpisodeService) GetDownloadLinks(animeName string, episodeNb int) (DownloadLinks, error) {
	params := url.Values{}
	payload := fmt.Sprintf(`n=%s\%d`, animeName, episodeNb)
	res, err := s.getEpisode(params, EpisodeDownloadPath, http.MethodPost, payload)
	if err != nil {
		return DownloadLinks{}, err
	}

	var dwnLinks DownloadLinks
	err = json.NewDecoder(res.Body).Decode(&dwnLinks)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	return dwnLinks, err
}
