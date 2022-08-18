package tohru

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"codeberg.org/omarkhatib/kobayashi"
	rncryptor "github.com/RNCryptor/RNCryptor-go"
)

const (
	GetEpisodePath      = "/anime/public/episodes/get-episodes-new"
	EpisodeDownloadPath = "/la/public/api/fw"
	BackupLinksPath     = "/anime/public/v-qs.php"
)

var (
	ErrBackupLink = fmt.Errorf("error while getting backup links")
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

func (s *EpisodeService) GetDirectDownloadLinks(animeName string, episodeNb int) (DownloadLinks, error) {
	return s.GetDirectDownloadLinksWithMax(animeName, episodeNb, -1)
}

func (s *EpisodeService) GetFirstDirectDownloadLink(animeName string, episodeNb int) (string, error) {
	link, err := s.GetDirectDownloadLinksWithMax(animeName, episodeNb, 1)
	if err != nil {
		return "", err
	}
	return link[0], err
}

type BackupLinks []struct {
	File  string `json:"file"`
	Label string `json:"label"`
}

func (s *EpisodeService) GetDirectDownloadLinksWithMax(animeName string, episodeNb int, maxNbOfLinks int) (DownloadLinks, error) {
	params := url.Values{}
	payload := fmt.Sprintf(`n=%s\%d`, animeName, episodeNb)
	res, err := s.getEpisode(params, EpisodeDownloadPath, http.MethodPost, payload)
	if err != nil {
		return DownloadLinks{}, err
	}

	var dwnLinks DownloadLinks
	err = json.NewDecoder(res.Body).Decode(&dwnLinks)
	if err != nil {
		return DownloadLinks{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	d := kobayashi.Decoder{}

	if len(dwnLinks) < maxNbOfLinks || maxNbOfLinks <= 0 {
		maxNbOfLinks = len(dwnLinks)
	}

	linksChan := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < len(dwnLinks); i++ {
		wg.Add(1)
		go func(link string) {
			url, _ := d.Decode(link)
			linksChan <- url
			wg.Done()
		}(dwnLinks[i])
	}

	go func() {
		wg.Wait()
		close(linksChan)
	}()

	var endRes DownloadLinks
	for link := range linksChan {
		if link == "" {
			continue
		}
		endRes = append(endRes, link)
		if len(endRes) == maxNbOfLinks {
			return endRes, nil
		}
	}

	if len(endRes) == 0 {
		return s.GetBackupLinks(animeName, episodeNb)
	}

	return endRes, nil
}

func (s *EpisodeService) GetBackupLinks(animeName string, episodeNb int) (DownloadLinks, error) {
	var endRes []string
	if s.client.cfg.backupLinksSecret != "" {
		apiUrl := BaseAPI
		resource := BackupLinksPath
		data := url.Values{}
		data.Set("f", animeName)
		data.Set("e", fmt.Sprintf("%d", episodeNb))
		data.Set("inf", `{"a": "UZDyJ8oRr3eFdNwJ0fynLrutZwyl89xCHORi2dy+/k5PtJyjg22p75JwZSyPtTSBdJCZo+WJiuR9gpxhTGWA0tZykydix3NY2UgfQPL9IXp3IL/VrNc2lS8rYUvXrCaUJEssegHOWS3+S8B3OFqSPg==", "b": "198.7.62.204"}`)

		u, _ := url.ParseRequestURI(apiUrl)
		u.Path = resource
		urlStr := u.String()

		client := &http.Client{}
		r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(r)
		if err != nil {
			return DownloadLinks{}, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return DownloadLinks{}, err
		}

		var backuplinks BackupLinks
		encrypted, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			return DownloadLinks{}, ErrBackupLink
		}
		decrypted, err := rncryptor.Decrypt(s.client.cfg.backupLinksSecret, encrypted)
		if err != nil {
			return DownloadLinks{}, ErrBackupLink
		}
		if err := json.Unmarshal(decrypted, &backuplinks); err != nil {
			return DownloadLinks{}, ErrBackupLink
		}
		for _, bl := range backuplinks {
			endRes = append(endRes, bl.File)
		}
	}
	if len(endRes) == 0 {
		return DownloadLinks{}, fmt.Errorf("all links are dead")
	}
	return endRes, nil
}
