# NPO Radio 2 Visual Radio stream link grabber

> This project is work in progress

This is a simple HTTP server that grabs the stream link from the NPO Radio 2 Visual Radio page and redirects to it. This because the stream link is not directly available.

## Usage

```bash
bun install
bun run index.ts

```

After that, the server is available on port 8080. You can get the stream link by requesting `/nporadio2.m3u8`.

## Environment variables

| Name | Description | Default |
| ---- | ----------- | ------- |
| `GRAB_TIMEOUT` | The timeout for the grabber in milliseconds | `4000` |
| `GRAB_PORT` | The port to listen on | `8080` |


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

> As this project is grabbing stuff from the NPO Radio 2 website, I have licensed this specific project under a non-commercial license. Please only use this for personal use and do not host this publicly.
This project is licensed under the CC BY-NC-SA 4.0 license. See [LICENSE](LICENSE) for more information.