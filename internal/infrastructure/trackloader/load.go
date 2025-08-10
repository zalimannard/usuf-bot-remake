package trackloader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"time"

	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/util"
)

type ytdlpJSON struct {
	Title      string  `json:"title"`
	Uploader   string  `json:"uploader"`
	Channel    string  `json:"channel"`
	Artist     string  `json:"artist"`
	Creator    string  `json:"creator"`
	MusicTrack string  `json:"track"`
	Duration   float64 `json:"duration"` // seconds
	Thumbnail  string  `json:"thumbnail"`
	Thumbnails []struct {
		URL string `json:"url"`
	} `json:"thumbnails"`
	// For playlists
	Entries []struct {
		Title      string `json:"title"`
		URL        string `json:"url"`
		WebpageURL string `json:"webpage_url"`
	} `json:"entries"`
}

func (l *Loader) Load(ctx context.Context, targetURL url.URL) (*track.Track, error) {
	// 1) Получаем JSON через yt-dlp
	data, effectiveURL, err := fetchYtDLP(ctx, targetURL.String())
	if err != nil {
		return nil, err
	}

	// 2) Извлекаем поля
	// Заголовок
	title := firstNonEmpty(
		data.Title,
		data.MusicTrack, // иногда yt-dlp отдаёт название трека в этом поле
		"Без названия",
	)

	// Автор
	var authorPtr *string
	if author := firstNonEmpty(data.Artist, data.Uploader, data.Channel, data.Creator); author != "" {
		authorPtr = util.Ptr(author)
	}

	// Длительность
	var durationPtr *time.Duration
	if data.Duration > 0 {
		d := time.Duration(data.Duration * float64(time.Second))
		durationPtr = &d
	}

	// Картинка
	var imageURLPtr *url.URL
	if thumb := pickThumbURL(data); thumb != "" {
		if u, err := url.Parse(thumb); err == nil {
			imageURLPtr = u
		}
	}

	// Итоговый URL трека
	finalURL := targetURL
	if effectiveURL != "" {
		if u, err := url.Parse(effectiveURL); err == nil {
			finalURL = *u
		}
	}

	// 3) Возвращаем доменную сущность
	return track.New(
		nil,
		title,
		authorPtr,
		durationPtr,
		finalURL,
		imageURLPtr,
	), nil
}

func fetchYtDLP(ctx context.Context, rawURL string) (ytdlpJSON, string, error) {
	// Запускаем yt-dlp -J чтобы получить метаданные
	data, err := runYtDLPJSON(ctx, rawURL)
	if err != nil {
		return ytdlpJSON{}, "", err
	}

	// Если это плейлист — берём первый элемент.
	if len(data.Entries) > 0 {
		entryURL := firstNonEmpty(data.Entries[0].WebpageURL, data.Entries[0].URL)
		if entryURL == "" {
			return ytdlpJSON{}, "", fmt.Errorf("yt-dlp: playlist entry without URL")
		}
		// Догружаем полные метаданные конкретного элемента
		item, err := runYtDLPJSON(ctx, entryURL)
		return item, entryURL, err
	}

	return data, rawURL, nil
}

func runYtDLPJSON(ctx context.Context, rawURL string) (ytdlpJSON, error) {
	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	cmd := exec.CommandContext(ctx, "yt-dlp", "-J", rawURL)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return ytdlpJSON{}, fmt.Errorf("yt-dlp failed: %w; stderr: %s", err, errBuf.String())
	}

	var data ytdlpJSON
	if err := json.Unmarshal(outBuf.Bytes(), &data); err != nil {
		return ytdlpJSON{}, fmt.Errorf("failed to parse yt-dlp JSON: %w", err)
	}

	return data, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func pickThumbURL(d ytdlpJSON) string {
	if d.Thumbnail != "" {
		return d.Thumbnail
	}
	// Часто последний элемент — самый крупный
	for i := len(d.Thumbnails) - 1; i >= 0; i-- {
		if d.Thumbnails[i].URL != "" {
			return d.Thumbnails[i].URL
		}
	}
	return ""
}
