package youtubeinforequester

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"usuf-bot-remake/internal/util"
)

func (r *Requester) GetURLByQuery(ctx context.Context, query string) (*url.URL, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("empty query")
	}
	limit := 1

	const endpoint = "https://www.googleapis.com/youtube/v3/search"

	params := url.Values{}
	params.Set("part", "id")       // нужен videoId
	params.Set("type", "video")    // ищем только видео
	params.Set("maxResults", "50") // максимум за один запрос
	params.Set("q", query)
	params.Set("key", r.apiKey)

	client := &http.Client{}
	collected := make([]url.URL, 0, limit)

	for {
		if len(collected) >= limit {
			break
		}

		reqURL := endpoint + "?" + params.Encode()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("youtube search request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			return nil, fmt.Errorf("youtube search returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
		}

		var payload struct {
			NextPageToken string `json:"nextPageToken"`
			Items         []struct {
				ID struct {
					VideoID string `json:"videoId"`
				} `json:"id"`
			} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return nil, fmt.Errorf("decode youtube search response: %w", err)
		}

		for _, it := range payload.Items {
			if it.ID.VideoID == "" {
				continue
			}
			u, err := url.Parse("https://www.youtube.com/watch?v=" + it.ID.VideoID)
			if err != nil {
				// маловероятно, но оборвёмся с понятной ошибкой
				return nil, fmt.Errorf("build video url for %q: %w", it.ID.VideoID, err)
			}
			collected = append(collected, *u)
			if len(collected) >= limit {
				break
			}
		}

		if payload.NextPageToken == "" {
			break
		}
		params.Set("pageToken", payload.NextPageToken)
	}

	// если ничего не нашли — вернём ошибку
	if len(collected) == 0 {
		return nil, fmt.Errorf("no results for query %q", query)
	}
	// подстрахуемся, чтобы не вернуть больше, чем просили
	if len(collected) > limit {
		collected = collected[:limit]
	}
	return util.Ptr(collected[0]), nil
}
