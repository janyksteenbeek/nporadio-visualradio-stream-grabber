# NPO Radio Visual Radio stream link grabber

This is a Go HTTP service that NPO Radio stream links with DRM. It saves the stream links into cache and updates them every 2 hours. This is done to prevent the stream links from expiring. 

## Usage

```bash
go mod download
go build -o nporadio-visualradio-stream-grabber cmd/grabber/main.go
./nporadio-visualradio-stream-grabber
```

After that, the server is available on port 8080.

## Streams

The following streams are available:

| URL                 | Station      | Format           | DRM Type            |
|---------------------|--------------|------------------|---------------------|
| `/nporadio1.m3u8`   | NPO Radio 1  | M3U8 stream      | FairPlay DRM (HLS)  |
| `/nporadio1.mpd`    | NPO Radio 1  | MPEG-DASH stream | Widevine DRM (DASH) |
| `/nporadio2.m3u8`   | NPO Radio 2  | M3U8 stream      | FairPlay DRM (HLS)  |
| `/nporadio2.mpd`    | NPO Radio 2  | MPEG-DASH stream | Widevine DRM (DASH) |
| `/npo3fm.m3u8`      | NPO 3FM      | M3U8 stream      | FairPlay DRM (HLS)  |
| `/npo3fm.mpd`       | NPO 3FM      | MPEG-DASH stream | Widevine DRM (DASH) |
| `/npoklassiek.m3u8` | NPO Klassiek | M3U8 stream      | FairPlay DRM (HLS)  |
| `/npoklassiek.mpd`  | NPO Klassiek | MPEG-DASH stream | Widevine DRM (DASH) |
| `/funx.m3u8`        | FunX         | M3U8 stream      | FairPlay DRM (HLS)  |
| `/funx.mpd`         | FunX         | MPEG-DASH stream | Widevine DRM (DASH) |



## Environment variables

| Name           | type            | Description                                 | Default            |
|----------------|-----------------|---------------------------------------------|--------------------|
| `GRAB_TIMEOUT` | duration string | The timeout for the grabber in milliseconds | ` 4 * time.Second` |
| `GRAB_PORT`    | int             | The port to listen on                       | `8080`             |
| `GRAB_REFRESH` | duration string | The refresh interval                        | `2 * time.Hour`    |


## Button for Home Assistant

I'm using this in Home Assistant to play easily start the stream on my Apple TV. This is the button I use:

```yaml
show_name: true
show_icon: false
type: button
name: ▶️ Radio 2
tap_action:
    action: call-service
    service: media_player.play_media
    target:
        device_id: <device_id>
    data:
        media_content_id: >-
            http://<your_ip>:8080/nporadio2.m3u8
        media_content_type: playlist
```

Make sure to replace `<device_id>` and `<your_ip>` with the correct values.

## License

> As this project is grabbing stuff from the NPO Radio websites, I have licensed this specific project under a non-commercial license. Please only use this for personal use and do not host this publicly.
This project is licensed under the CC BY-NC-SA 4.0 license. See [LICENSE](LICENSE) for more information.