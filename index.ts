import { createServer, IncomingMessage, ServerResponse } from 'http';
import puppeteer, {Browser} from "puppeteer";

const TIMEOUT: number = parseInt(process.env.GRAB_TIMEOUT || '') || (4 * 1000);
const PORT: number = parseInt(process.env.GRAB_PORT || '') || 8080;

interface StreamLinkRequest {
    profileName: string;
    drmType: string;
    referrerUrl: string;
}

interface StreamLinkResponse {
    stream: {
        streamURL: string;
    };
}

interface RouteConfig {
    drmType: string;
    profileName: string;
}

const routeConfigs: { [key: string]: RouteConfig } = {
    '/nporadio2.m3u8': { drmType: 'fairplay', profileName: 'hls' },
    '/nporadio2.mpd': { drmType: 'widevine', profileName: 'dash' }
};

async function fetchStreamUrl(browser: Browser, config: RouteConfig): Promise<string | null> {
    const page = await browser.newPage();
    await page.setRequestInterception(true);

    page.on('request', (interceptedRequest) => {
        if (interceptedRequest.url().includes('https://prod.npoplayer.nl/stream-link') && interceptedRequest.method() === 'POST') {
            const postData: StreamLinkRequest = {
                ...JSON.parse(interceptedRequest.postData() || '{}'),
                drmType: config.drmType,
                profileName: config.profileName
            };
            interceptedRequest.continue({
                method: 'POST',
                postData: JSON.stringify(postData),
                headers: interceptedRequest.headers()
            });
        } else {
            interceptedRequest.continue();
        }
    });

    let streamUrl: string | null = null;
    page.on('response', async (response) => {
        if (response.request().method() === 'OPTIONS') return;
        if (response.url().includes('https://prod.npoplayer.nl/stream-link')) {
            const jsonResponse: StreamLinkResponse = await response.json();
            streamUrl = jsonResponse.stream.streamURL;
        }
    });

    await page.goto('https://www.nporadio2.nl/live', { waitUntil: 'networkidle2' });
    await new Promise((resolve) => setTimeout(resolve, TIMEOUT));
    await page.close();

    return streamUrl;
}

createServer(async (req, res) => {
    const browser = await puppeteer.launch({ headless: 'new' });
    const date= new Date().toISOString();

    const config = routeConfigs[req.url || ''];

    if (config === undefined) {
        console.log(`[${date}] 404 ${req.url}`);
        res.writeHead(404).end('Not Found');
    }
    const streamUrl = await fetchStreamUrl(browser, config);
    if (streamUrl) {
        console.log(`[${date}] 302 ${req.url} -> ${streamUrl}`);
        res.writeHead(302, { 'Location': streamUrl }).end();
    } else {
        console.log(`[${date}][${date}] 408 ${req.url}`);
        res.writeHead(408).end('Request Timeout');
        res.end('Stream URL not found in time');
    }

}).listen(PORT, () => {
    console.log(`Server running at port ${PORT} 🚀`);
});
