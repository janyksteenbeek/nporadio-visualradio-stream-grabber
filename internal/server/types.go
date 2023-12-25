package server

import "time"

type StreamLinkRequest struct {
	ProfileName string `json:"profileName"`
	DrmType     string `json:"drmType"`
	ReferrerUrl string `json:"referrerUrl"`
}

type StreamLinkResponse struct {
	Stream struct {
		StreamURL string `json:"streamURL"`
	} `json:"stream"`
}

type RouteConfig struct {
	DrmType          string
	ProfileName      string
	StreamBuilderURL string
	URL              string
}

type NextData struct {
	Props struct {
		PageProps struct {
			Player struct {
				TokenURL string `json:"tokenUrl"`
			} `json:"player"`
		} `json:"pageProps"`
	} `json:"props"`
}

type PlayerTokenResponse struct {
	PlayerToken string `json:"playerToken"`
}

type CachedStreamUrl struct {
	StreamURL string
	Timestamp time.Time
}
