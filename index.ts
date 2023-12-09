import * as http from 'http';
import * as puppeteer from 'puppeteer';

const TIMEOUT: number = parseInt(process.env.GRAB_TIMEOUT || '') || (4 * 1000);
const PORT: number = parseInt(process.env.GRAB_PORT || '') || 8080;

http.createServer(async (req: http.IncomingMessage, res: http.ServerResponse): Promise<void> => {
    if (req.url === '/nporadio2.m3u8') {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        let receivedStreamUrl: boolean = false;

        page.on('response', async (response: puppeteer.HTTPResponse) => {
            // Check if preflight
            if (response.status() === 204) {
                return;
            }
            
            const url: string = response.url();
            if (url.includes('https://prod.npoplayer.nl/stream-link')) {
                const jsonResponse: any = await response.json();
                const streamURL: string = jsonResponse?.stream?.streamURL;

                if (streamURL) {
                    res.writeHead(302, { 'Location': streamURL });
                    res.end();
                    receivedStreamUrl = true;
                }
            }
        });

        await page.goto('https://www.nporadio2.nl/live', {
            waitUntil: 'networkidle2'
        });

        setTimeout(() => {
            if (!receivedStreamUrl) {
                res.writeHead(408);
                res.end('Stream URL not found in time');
            }
            browser.close();
        }, TIMEOUT);

    } else {
        res.writeHead(404);
        res.end('Not Found');
    }
}).listen(PORT, () => {
    console.log(`Server running at port ${PORT} 🚀`);
    console.log('To grab the stream URL, visit http://<ip-address>:8080/nporadio2.m3u8');
    console.log('To quit, press CTRL+C');
});
