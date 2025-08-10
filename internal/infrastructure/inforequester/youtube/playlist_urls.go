package youtubeinforequester

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (r *Requester) PlaylistURLs(ctx context.Context, playlistID string) ([]url.URL, error) {
	const endpoint = "https://www.googleapis.com/youtube/v3/playlistItems"

	params := url.Values{}
	params.Set("part", "contentDetails") // нужен videoId
	params.Set("maxResults", "50")       // максимум на страницу
	params.Set("playlistId", playlistID)
	params.Set("key", r.apiKey)

	client := &http.Client{}
	result := make([]url.URL, 0, 64)

	for {
		reqURL := endpoint + "?" + params.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("youtube api request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			return nil, fmt.Errorf("youtube api returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
		}

		var payload struct {
			NextPageToken string `json:"nextPageToken"`
			Items         []struct {
				ContentDetails struct {
					VideoID string `json:"videoId"`
				} `json:"contentDetails"`
			} `json:"items"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return nil, fmt.Errorf("decode youtube response: %w", err)
		}

		for _, it := range payload.Items {
			vid := strings.TrimSpace(it.ContentDetails.VideoID)
			if vid == "" {
				continue
			}
			u, err := url.Parse("https://www.youtube.com/watch?v=" + vid)
			if err != nil {
				return nil, fmt.Errorf("build video url for %q: %w", vid, err)
			}
			result = append(result, *u)
		}

		if payload.NextPageToken == "" {
			break
		}
		params.Set("pageToken", payload.NextPageToken)
	}

	return result, nil
}
