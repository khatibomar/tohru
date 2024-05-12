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

	"github.com/khatibomar/kobayashi"
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

type DownloadInfo struct {
	EpisodeHostLink           string
	EpisodeDirectDownloadLink string
}

type DownloadLinks []string

type DownloadInfos []DownloadInfo

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
	var payload JsonPayload
	var err error
	var payloadStr string

	err = payload.WithJustInfo("No")
	if err != nil {
		return []Episode{}, err
	}
	err = payload.WithAnimeId(animeID)
	if err != nil {
		return []Episode{}, err
	}
	payloadStr, err = payload.ToJson()
	if err != nil {
		return []Episode{}, err
	}
	res, err := s.getEpisode(params, GetEpisodePath, http.MethodPost, payloadStr)
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
	var payload JsonPayload
	var err error
	var payloadStr string

	err = payload.WithEpisodeId(episodeID)
	if err != nil {
		return Episode{}, err
	}
	err = payload.WithAnimeId(animeID)
	if err != nil {
		return Episode{}, err
	}
	payloadStr, err = payload.ToJson()
	if err != nil {
		return Episode{}, err
	}
	res, err := s.getEpisode(params, GetEpisodePath, http.MethodPost, payloadStr)
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

	var payload JsonPayload
	var err error
	var payloadStr string

	err = payload.WithN(animeName, episodeNb)
	if err != nil {
		return DownloadLinks{}, err
	}

	payloadStr, err = payload.ToJson()
	if err != nil {
		return DownloadLinks{}, err
	}

	res, err := s.getEpisode(params, EpisodeDownloadPath, http.MethodPost, payloadStr)
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

func (s *EpisodeService) GetDirectDownloadInfos(animeName string, episodeNb int) (DownloadInfos, error) {
	return s.GetDirectDownloadInfosWithMax(animeName, episodeNb, -1)
}

func (s *EpisodeService) GetFirstDirectDownloadInfo(animeName string, episodeNb int) (DownloadInfo, error) {
	link, err := s.GetDirectDownloadInfosWithMax(animeName, episodeNb, 1)
	if err != nil {
		return DownloadInfo{}, err
	}
	return link[0], err
}

type BackupLinks []struct {
	File  string `json:"file"`
	Label string `json:"label"`
}

func (s *EpisodeService) GetDirectDownloadInfosWithMax(animeName string, episodeNb int, maxNbOfLinks int) (DownloadInfos, error) {
	params := url.Values{}
	var payload JsonPayload
	var err error
	var payloadStr string

	err = payload.WithN(animeName, episodeNb)
	if err != nil {
		return DownloadInfos{}, err
	}

	payloadStr, err = payload.ToJson()
	if err != nil {
		return DownloadInfos{}, err
	}

	res, err := s.getEpisode(params, EpisodeDownloadPath, http.MethodPost, payloadStr)
	if err != nil {
		return DownloadInfos{}, err
	}

	var dwnLinks DownloadLinks
	err = json.NewDecoder(res.Body).Decode(&dwnLinks)
	if err != nil {
		return DownloadInfos{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	d := kobayashi.Decoder{}

	if len(dwnLinks) < maxNbOfLinks || maxNbOfLinks <= 0 {
		maxNbOfLinks = len(dwnLinks)
	}

	linksChan := make(chan DownloadInfo)
	var wg sync.WaitGroup

	for i := 0; i < len(dwnLinks); i++ {
		wg.Add(1)
		go func(link string) {
			url, _ := d.Decode(link)
			linksChan <- DownloadInfo{link, url}
			wg.Done()
		}(dwnLinks[i])
	}

	go func() {
		wg.Wait()
		close(linksChan)
	}()

	var endRes DownloadInfos
	for link := range linksChan {
		if link.EpisodeDirectDownloadLink == "" {
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

func (s *EpisodeService) GetBackupLinks(animeName string, episodeNb int) (DownloadInfos, error) {
	var endRes []DownloadInfo
	if s.client.cfg.backupLinksSecret != "" {
		apiUrl := BaseAPI
		resource := BackupLinksPath
		data := url.Values{}
		data.Set("f", animeName)
		data.Set("e", fmt.Sprintf("%d", episodeNb))
		data.Set("inf", `{"a": "4+mwbwVfA5wLr7a4GBQvzMy1/jO9fRQ/lKJXNS4vbW/FqNL3j0vtOPd5pQx2UxrJ/8UF0Xr/v/dxkse3tjvEg/1uLKKZM8CALrQrGtw0pQqZ+UiyBJqVXe9tlbFSkV9XQRkIC6qjY66uzkzk6wauPw==", "b": "217.138.207.148"}`)

		u, _ := url.ParseRequestURI(apiUrl)
		u.Path = resource
		urlStr := u.String()

		client := &http.Client{}
		r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(r)
		if err != nil {
			return DownloadInfos{}, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return DownloadInfos{}, err
		}

		var backuplinks BackupLinks
		encrypted, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			return DownloadInfos{}, ErrBackupLink
		}
		decrypted, err := rncryptor.Decrypt(s.client.cfg.backupLinksSecret, encrypted)
		if err != nil {
			return DownloadInfos{}, ErrBackupLink
		}
		if err := json.Unmarshal(decrypted, &backuplinks); err != nil {
			return DownloadInfos{}, ErrBackupLink
		}
		for _, bl := range backuplinks {
			endRes = append(endRes, DownloadInfo{"Backup link", bl.File})
		}
	}
	if len(endRes) == 0 {
		return DownloadInfos{}, fmt.Errorf("all links are dead")
	}
	return endRes, nil
}
