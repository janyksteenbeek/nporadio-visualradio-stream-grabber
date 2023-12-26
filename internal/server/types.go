package server

import "time"

// StreamLinkRequest contains the request to the NPO stream link builder
type StreamLinkRequest struct {
	ProfileName string `json:"profileName"`
	DrmType     string `json:"drmType"`
	ReferrerUrl string `json:"referrerUrl"`
}

// StreamLinkResponse contains the response with the stream URL from the NPO stream link builder
type StreamLinkResponse struct {
	Stream struct {
		StreamURL string `json:"streamURL"`
	} `json:"stream"`
}

// RouteConfig contains the configuration for a route, including the DRM type, profile name, stream builder URL and the URL to fetch the stream URL for
type RouteConfig struct {
	DrmType          string
	ProfileName      string
	StreamBuilderURL string
	URL              string
}

// NextData contains the NextJS page data for the NPO visual radio player page
type NextData struct {
	Props struct {
		PageProps struct {
			Player struct {
				TokenURL string `json:"tokenUrl"`
			} `json:"player"`
		} `json:"pageProps"`
	} `json:"props"`
}

// PlayerTokenResponse contains the response with the player token from the NPO player token service
type PlayerTokenResponse struct {
	PlayerToken string `json:"playerToken"`
}

// CachedStreamUrl contains a cached stream URL and the timestamp of the cache
type CachedStreamUrl struct {
	StreamURL string
	Timestamp time.Time
}
