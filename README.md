# The Midnight Lyrics bot
This bot will post lyrics from The Midnight songs on different social media platform.

The following platforms are supported:
- Bluesky
- Instagram
- Threads

Special thanks to [Carina](https://github.com/carina092) for the design of the images

Based on [Julia](https://github.com/juliajungle)'s original Node.js version: https://github.com/juliajungle/midnight-lyrics

## Bluesky

Command line flag: `-bluesky`

Environment variables required:
- `BOTSKY_HANDLE`
- `BOTSKY_APPKEY`

## Instagram

Command line flag: `-instagram`

Environment variables required:
- `INSTAGRAM_ACCESS_TOKEN`
- `INSTAGRAM_IMAGES_URL`

All images are saved in `generated_images` and a webserver serving these images is required.
The `INSTAGRAM_IMAGES_URL` needs to be a publicly accessible URL on the internet, so Instagram can fetch the images and post them.

## Threads

Command line flag: `-threads`

Environment variables required:
- `THREADS_ACCESS_TOKEN`

## Other command line flags

- `-generate-all-images` - Used to generate all the images needed for the Instagram part of the bot. Only useful for testing, as the bot will generate images as needed.
- `-fake-midnight` - Fake that the time is 00:00. Only useful for testing of the special post at midnight.
