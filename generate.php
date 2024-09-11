<?php
require './vendor/autoload.php';
require 'inc/functions.php';

$albums = [
	'days_of_thunder',
	'endless_summer',
	'heroes',
	'horror_show',
	'kids',
	'monsters',
	'nocturnal',
	'songs', //Special album for songs that do not fit under a real album.
];

if (!empty($_SERVER['argv'][1]) && $_SERVER['argv'][1] == 'generate_all') {
	foreach ($albums as $album_name) {
		$album = json_decode(file_get_contents('albums/'.$album_name.'.json'));

		foreach ($album->songs as $song) {
			foreach (explode('|', $song->lyrics) as $lyrics) {
				generate_album_image($album_name, $album, $song, $lyrics);
			}
		}
	}
}

if (!empty($_SERVER['argv'][1]) && !empty($_SERVER['argv'][2]) && $_SERVER['argv'][1] == 'generate_album') {
	$album_name = $_SERVER['argv'][2];

	if (!in_array($album_name, $albums))
		exit($album_name.' is not a valid album.');

	$album = json_decode(file_get_contents('albums/'.$album_name.'.json'));

	foreach ($album->songs as $song) {
		foreach (explode('|', $song->lyrics) as $lyrics) {
			generate_album_image($album_name, $album, $song, $lyrics);
		}
	}
}

else {
	$random_album_name = $albums[random_int(0, count($albums) - 1)];
	$random_album = json_decode(file_get_contents('albums/'.$random_album_name.'.json'));

	$random_song = $random_album->songs[random_int(0, count($random_album->songs) - 1)];
	$random_song_lyrics_split = explode('|', $random_song->lyrics);
	$random_lyrics = $random_song_lyrics_split[random_int(0, count($random_song_lyrics_split) - 1)];

	generate_album_image($random_album_name, $random_album, $random_song, $random_lyrics);
}
?>
