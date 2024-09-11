<?php
use Intervention\Image\Drivers\Gd\Driver;
use Intervention\Image\ImageManager;
use Intervention\Image\Typography\FontFactory;

function generate_album_image($album_name, $album, $song, $lyrics) {
	$manager = new ImageManager(new Driver());
	$image = $manager->read('images/'.$album_name.'.png');

	$image->text($lyrics, 540, 540, function (FontFactory $font){
		$font->filename('fonts/Optiker-K.ttf');
		$font->color('F4F4F4');
		$font->size(52);
		$font->align('center');
		$font->valign('middle');
		$font->lineHeight(3);
		$font->stroke('A61F59', 1);
		$font->wrap(1080 - (78 * 2));
	});

	//They are not a "real" album, so don't show it.
	if ($album_name == 'songs')
		$info_text = $song->title;
	else
		$info_text = $album->title."\n".$song->title;

	$image->text($info_text, 78, 952, function (FontFactory $font){
		$font->filename('fonts/Optiker-K.ttf');
		$font->color('F4F4F4');
		$font->size(40);
		$font->valign('middle');
		$font->lineHeight(2.8);
		$font->stroke('A61F59', 1);
		$font->wrap(708);
	});

	$folder = 'generated_images/'.$album_name.'/'.strtolower(str_replace(' ', '_', $song->title));

	if (!is_dir($folder))
		mkdir($folder, 0777, true);

	$image->toPng()->save($folder.'/'.sha1($lyrics).'.png');
}
?>
