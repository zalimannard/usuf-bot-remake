package trackprovider

import (
	"context"
	"net/url"
)

func (p *Provider) GetURLByQuery(ctx context.Context, query string) (*url.URL, error) {
	return p.infoRequester.GetURLByQuery(ctx, query)
}
