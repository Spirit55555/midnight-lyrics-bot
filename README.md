# The Midnight Lyrics bot
This bot will post lyrics from The Midnight songs, on different social media platform. Currently on Bluesky and Threads.

# How to set it up
Download the artifact `lyrics-bot` under actions, and run it on Linux.\
A cronjob is currently the way to make this run, it needs to be setup like this.\
This is what needs to be put in the crontab file using `crontab -e`.\
Remember to make the file `midnight-lyrics-bot` executable `chmod +x midnight-lyrics-bot`.

The `INSTAGRAM_IMAGES_URL` needs to be a publicly accessible URL on the internet, for the API to acccess and post that image.
```
INSTAGRAM_ACCESS_TOKEN=
INSTAGRAM_IMAGES_URL=
BOTSKY_HANDLE=
BOTSKY_APPKEY=
THREADS_ACCESS_TOKEN=
0 */4 * * * cd /home/USERNAME/midnight-lyrics-bot && ./midnight-lyrics-bot
```
## Sepcify services to make the bot run on
You need to specify what services to run it on, this can be done using the following arguments after running the binary. This is done to seperate the times the bot posts.
- -instagram
- -bluesky
- -threads


## Arguments
These can be used when running the binary
- -generate-all-images - Used to generate all the images needed for the Instagram part of the bot
