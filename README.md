# The Midnight Lyrics bot
This bot will post song lyrics from [The Midnight](https://themidnightofficial.com/) on different social media platform.

The following platforms are currently supported:
- Bluesky
- Instagram
- Threads

Special thanks to [Carina](https://github.com/carina092) for the design of the images.
Based on [Julia](https://github.com/juliajungle)'s original [Node.js version](https://github.com/juliajungle/midnight-lyrics).

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

If `INSTAGRAM_IMAGES_URL` is set to `https://instagram.example.com` then the image `generated_images/album/song/sha1hash.jpg` should be accessible on `https://instagram.example.com/album/song/sha1hash.jpg`

## Threads

Command line flag: `-threads`

Environment variables required:
- `THREADS_ACCESS_TOKEN`

## Other command line flags

- `-refresh-tokens` - USe to refresh and print new tokens for Threads and Instagram.
- `-generate-image` - Used to generate the image for the selected (or random) lyric. Album and song number can be specified. Only useful for testing.
- `-generate-all-images` - Used to generate all the images needed for the Instagram part of the bot. Only useful for testing, as the bot will generate images as needed.
- `-fake-midnight` - Fake that the time is 00:00. Only useful for testing of the special post at midnight.

## Specify album and song

If you want a specific album and (optionally) song selected when running the bot, it can take them as parameters like so:

```
midnight-lyrics-bot days_of_thunder 6
```

This will select a random lyric from the song Los Angeles from the Days of Thunder album.

## Docker

Assuming you are in the source directory, the bot can be built and installed via [Docker](https://www.docker.com/) the following way:

```
docker buildx build -t midnight-lyrics-bot .
```

### Running the bot

Running the bot requires the above mentioned environment variables set, which can be done via `--env` or (preferably) via `--env-file`:

```
docker run --env-file=socials.env midnight-lyrics-bot -bluesky
```

Example `socials.env` file for Bluesky:

```
BOTSKY_HANDLE=myhandle@example.com
BOTSKY_APPKEY=P@passw0rd
```

#### Instagram images

Running it for Instagram requires mounting of the `generated_images` directory to a directory where the images are served by a webserver on the host machine ([see note above](#instagram)):

```
docker run --env-file=socials.env -v /opt/midnight-lyrics-bot/images:/generated_images midnight-lyrics-bot -instagram
```

Replace `/opt/midnight-lyrics-bot/images` with the path on the host machine.

#### Special post at midnight

In order to get the special post at midnight (00:00) at the expected time, the timezone has to be set to the the expected one.

The following line can be added to your `socials.env` file:

```
TZ=Europe/Copenhagen
```

A list of valid timezones can be found [here](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).
