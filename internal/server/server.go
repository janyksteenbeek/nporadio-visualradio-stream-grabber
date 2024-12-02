package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-rod/rod"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultTimeout         = 4 * time.Second
	defaultPort            = 8080
	defaultRefreshInterval = 2 * time.Hour
)

var routeConfigs = map[string]RouteConfig{
	"/nporadio1.m3u8":   {DrmType: "fairplay", ProfileName: "hls", URL: "https://www.nporadio1.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
	"/nporadio1.mpd":    {DrmType: "widevine", ProfileName: "dash", URL: "https://www.nporadio1.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
	"/nporadio2.m3u8":   {DrmType: "fairplay", ProfileName: "hls", URL: "https://www.nporadio2.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
	"/nporadio2.mpd":    {DrmType: "widevine", ProfileName: "dash", URL: "https://www.nporadio2.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
	"/npo3fm.m3u8":      {DrmType: "fairplay", ProfileName: "hls", URL: "https://www.npo3fm.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
	"/npo3fm.mpd":       {DrmType: "widevine", ProfileName: "dash", URL: "https://www.npo3fm.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
	"/npoklassiek.m3u8": {DrmType: "fairplay", ProfileName: "hls", URL: "https://www.nporadio4.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
	"/npoklassiek.mpd":  {DrmType: "widevine", ProfileName: "dash", URL: "https://www.nporadio4.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
	"/funx.m3u8":        {DrmType: "fairplay", ProfileName: "hls", URL: "https://www.funx.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
	"/funx.mpd":         {DrmType: "widevine", ProfileName: "dash", URL: "https://www.funx.nl/live", StreamBuilderURL: "https://prod.npoplayer.nl/stream-link"},
}

var cachedUrls = make(map[string]CachedStreamUrl)
var port int
var refreshInterval time.Duration
var timeout time.Duration
var mutex = &sync.Mutex{}

func fetchStreamUrl(ctx context.Context, config RouteConfig, authToken string) (string, error) {

	client := &http.Client{
		Timeout: timeout,
	}

	requestBody, err := json.Marshal(StreamLinkRequest{
		ProfileName: config.ProfileName,
		DrmType:     config.DrmType,
		ReferrerUrl: config.URL,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.StreamBuilderURL, strings.NewReader(string(requestBody)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch stream URL, status code: " + resp.Status)
	}

	var response StreamLinkResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	log.Printf("Fetched stream URL for %s (%s): %s", config.URL, config.ProfileName, response.Stream.StreamURL)

	return response.Stream.StreamURL, nil
}

func fetchPlayerToken(ctx context.Context, tokenURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", tokenURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch player token, status code: %d", resp.StatusCode)
	}

	var tokenResponse PlayerTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.PlayerToken, nil
}

func fetchTokenURL(ctx context.Context, browser *rod.Browser, livePageURL string) (string, error) {
	page := browser.MustPage(livePageURL)
	defer page.Close()
	page.MustWaitLoad()

	scriptContent := page.MustElement("#__NEXT_DATA__").MustText()
	var nextData NextData
	err := json.Unmarshal([]byte(scriptContent), &nextData)
	if err != nil {
		return "", err
	}

	return nextData.Props.PageProps.Player.TokenURL, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s (%s, %s)", r.Method, r.URL.Path, r.RemoteAddr, r.Header.Get("User-Agent"))

	mutex.Lock()
	cached, ok := cachedUrls[r.URL.Path]
	mutex.Unlock()

	if !ok || time.Since(cached.Timestamp) >= refreshInterval {
		http.Error(w, "Stream URL not found or expired", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, cached.StreamURL, http.StatusFound)
}

func updateStreamUrls(ctx context.Context, browser *rod.Browser) {

	for path, config := range routeConfigs {
		tokenUrl, err := fetchTokenURL(ctx, browser, config.URL)
		if err != nil {
			log.Printf("Error fetching token URL for %s (%s, %s): %v", path, config.URL, config.ProfileName, err)
			continue
		}

		playerToken, err := fetchPlayerToken(ctx, tokenUrl)
		if err != nil {
			log.Printf("Error fetching player token for %s (%s, %s): %v", path, config.URL, config.ProfileName, err)
			continue
		}

		streamUrl, err := fetchStreamUrl(ctx, config, playerToken)
		if err != nil {
			log.Printf("Error fetching stream URL for %s (%s, %s): %v", path, config.URL, config.ProfileName, err)
			continue
		}

		mutex.Lock()
		cachedUrls[path] = CachedStreamUrl{
			StreamURL: streamUrl,
			Timestamp: time.Now(),
		}
		mutex.Unlock()
	}
}

func startUpdateTicker(ctx context.Context, browser *rod.Browser) {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateStreamUrls(ctx, browser)
		}
	}
}

func StartServer() {
	port, _ = strconv.Atoi(os.Getenv("PORT"))
	timeout, _ = time.ParseDuration(os.Getenv("TIMEOUT"))
	refreshInterval, _ = time.ParseDuration(os.Getenv("REFRESH_INTERVAL"))

	if port == 0 {
		port = defaultPort
	}
	if timeout == 0 {
		timeout = defaultTimeout
	}
	if refreshInterval == 0 {
		refreshInterval = defaultRefreshInterval
	}
	browser := rod.New()

	err := browser.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to browser: %v", err)
	}
	defer browser.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go updateStreamUrls(ctx, browser)
	go startUpdateTicker(ctx, browser)

	http.HandleFunc("/", handleRequest)
	log.Printf("Server running at port %d 🚀", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
