# The Midnight Lyrics bot
This bot will post lyrics from The Midnight songs, on different social media platform. Currently on Bluesky and Threads.

# How to set it up
Download the artifact under actions, and run it on Linux.\
A cronjob is currently the way to make this run, it needs to be setup like this.\
This is what needs to be put in the crontab file using `crontab -e`.\
Remember to make the file `midnight-lyrics-bot` executable `chmod +x midnight-lyrics-bot`.
```
BOTSKY_HANDLE=
BOTSKY_APPKEY=
THREADS_ACCESS_TOKEN=
0 */4 * * * cd /home/jens/midnight-lyrics-bot && ./midnight-lyrics-bot
```
