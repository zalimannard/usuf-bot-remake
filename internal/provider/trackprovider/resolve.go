package trackprovider

import (
	"context"
	"net/url"
	"strings"
)

func (p *Provider) ExpandURL(ctx context.Context, targetURL url.URL) ([]url.URL, error) {
	host := strings.ToLower(targetURL.Host)

	switch {
	case isYouTube(host):
		if listID := targetURL.Query().Get("list"); listID != "" {
			return p.infoRequester.PlaylistURLs(ctx, listID)
		}
		return []url.URL{targetURL}, nil
	default:
		return []url.URL{targetURL}, nil
	}

}

func isYouTube(host string) bool {
	return strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be")
}
