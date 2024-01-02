# Instructions
Install all packages `make install`

Start server `make server`


# TODO
1. Add swagger docs
2. Validate user input
2. Add auth so that not any user can generate the video
3. OWASP top 10 vulnerabilities. e.g. API rate limiting
4. Gather performance requirements. (rps, latency, error rate...) - load testing

# Commands
To generate the video using plain cli: `ffmpeg -loop 1 -i cover.jpg -i 1.jpg -i 2.jpg -i Jingle-Bells.mp3 -c:v libx264 -tune stillimage -c:a aac -b:a 192k -pix_fmt yuv420p -shortest exported.mp4`

with timeout version on linux with glob`ffmpeg -framerate 1 -pattern_type glob -i '*.jpg' -i Jingle-Bells.mp3 -c:v libx264 -t 15 -pix_fmt yuv420p -vf scale=320:240 out.mp4`

with timeout version on Windows ends when slids finishes `ffmpeg -y -framerate 1 -i 'upload/%d.jpg' -i Jingle-Bells.mp3 -c:v libx264 -pix_fmt yuv420p -vf scale=320:240 -t 15 -shortest D:/videos/upload/out.mp4`

Fade effects for videos or images 
To generate the command from UI, use https://ffmpeg.shanewhite.co/
```bash
ffmpeg -i 4.mp4 \
-filter_complex "[0:v]drawtext=fontfile=Lato-Light.ttf:text='Hello world':fontsize=130:fontcolor=ffffff:alpha='if(lt(t,1),0,if(lt(t,2),(t-1)/1,if(lt(t,32),1,if(lt(t,33),(1-(t-32))/1,0))))':x=(w-text_w)/2:y=(h-text_h)/2" out.mp4
```


Concatenate the vidoes and images together in upload-1 with AUDIO:
```bash
ffmpeg -y \
-loop 1 -framerate 24 -t 10 -i 1.jpg \
-i 4.mp4 \
-loop 1 -framerate 24 -t 10 -i 2.jpg \
-loop 1 -framerate 24 -t 10 -i 3.jpg \
-f lavfi -t 0.1 -i anullsrc=channel_layout=stereo:sample_rate=44100 \
-filter_complex "[0:v][4:a][1:v][1:a][2:v][4:a][3:v][4:a]concat=n=4:v=1:a=1" out.mp4
```

Without audio
```bash
ffmpeg -y \
-loop 1 -framerate 24 -t 10 -i 1.jpg \
-i 4.mp4 \
-loop 1 -framerate 24 -t 10 -i 2.jpg \
-loop 1 -framerate 24 -t 10 -i 3.jpg \
-filter_complex "[0][1][2][3]concat=n=4:v=1:a=0" out.mp4
```

Superuser by setting same sar value with background audio (working 1)
while 
- keep the image/video aspect ratio 
- padding if there is any blank spaces.
```bash
ffmpeg -y \
-loop 1 -framerate 24 -t 10 -i 1.jpg \
-i 4.mp4 \
-loop 1 -framerate 24 -t 10 -i 2.jpg \
-i Jingle-Bells.mp3 \
-filter_complex "[0]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[im];[1:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[vid];[2]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[im1];[im][vid][im1]concat=n=3:v=1:a=0" -shortest out.mp4

# v2 with audio merged and mapped rather than concat with audio directly, adding 4.mp4
ffmpeg -y \
-loop 1 -framerate 24 -t 10 -i 1.jpg \
-i 4.mp4 \
-loop 1 -framerate 24 -t 10 -i 2.jpg \
-i Jingle-Bells.mp3 \
-i 6.mp4 \
-filter_complex "[0]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[im];[1:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[vid];[2:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[im1];[4:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[vid1];[im][vid][im1][vid1]concat=n=4:v=1:a=0[v];[3:a]amerge=inputs=1[a]" \
-map "[v]" -map "[a]" -ac 2 output.mp4


# v3: v2 + cover page with fade in effect for titles 
# Will need to know the 
ffmpeg -y \
-loop 1 -framerate 24 -t 10 -i 1.jpg \
-i 4.mp4 \
-loop 1 -framerate 24 -t 10 -i 2.jpg \
-i /Users/yangliu/vidoes/media/theme/ocean/bgm.mp3 \
-i 6.mp4 \
-i cover.mp4 \
-filter_complex "[0]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[im];[1:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[vid];[2:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[im1];[4:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[vid1];[5:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[cover];[cover][im][vid][im1][vid1]concat=n=5:v=1:a=0[v];[3:a]amerge=inputs=1[a]" \
-map "[v]" -map "[a]" -ac 2 output.mp4

# v4: v3 on Mac quickplayer with shortest first
ffmpeg -y \
-framerate 0.5 \
-pix_fmt yuv420p \
-loop 1 -framerate 24 -t 10 -i 1.jpg \
-i 4.mp4 \
-loop 1 -framerate 24 -t 10 -i 2.jpg \
-i /Users/yangliu/vidoes/media/theme/ocean/bgm.mp3 \
-i 6.mp4 \
-i cover.mp4 \
-filter_complex "[0]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[im];[1:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[vid];[2:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[im1];[4:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[vid1];[5:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1[cover];[cover][im][vid][im1][vid1]concat=n=5:v=1:a=0[v];[3:a]amerge=inputs=1[a]" \
-map "[v]" -map "[a]" -ac 2 -shortest output.mp4
```

Fade in & out effect
```bash
# Experiment
ffmpeg -y \
-i "6.mp4" -loop 1 \
-i "1.jpg" \
-filter_complex "[1:v]scale=w='if(between(t,10,14),1469-(1469-293)*(t-10)/4,if(lt(t,10),1469,293))':h='if(between(t,10,14),856-(856-171)*(t-10)/4,if(lt(t,10),856,171))':eval=frame[img];[0:v][img]overlay=x='if(between(t,10,14),(W-w)/2-((W-w)/2-10)*(t-10)/4,if(lt(t,10),(W-w)/2,10))':y='if(between(t,10,14),(H-h)/2-((H-h)/2-(H-h-40))*(t-10)/4,if(lt(t,10),(H-h)/2,H-h-40))':shortest=1" \
-crf 18 -c:a copy out.mp4
```


```bash
ffmpeg -y \
-loop 1 -framerate 24 -t 10 -i 1.jpg \
-i 4.mp4 \
-i Jingle-Bells.mp3 \
-filter_complex "[0]scale=432:432,setsar=1[im];[1:v]scale=432:432,setsar=1[vid];[2:a]asplit=1[aud];[im][vid][aud]concat=n=3:v=1:a=2" out.mp4
```

Superuser by setting same sar value using map option
```bash
ffmpeg -y \
-loop 1 -framerate 24 -t 10 -i 1.jpg \
-i 4.mp4 \
-i Jingle-Bells.mp3 \
-filter_complex "[0]scale=432:432,setsar=1[im];[1:v]scale=432:432,setsar=1[vid];[2:a]asplit=1[aud]" \
-map [im] -map [vid] -map [aud] out.mp4
```

# Useful sites
[Sample videos](https://www.pexels.com/search/videos/ocean/)
[mp3 download](https://mp3juice.ms/)
[Free stock image download](https://unsplash.com/s/photos/friend)
[Stock jpg image donwload](https://www.istockphoto.com/search/2/image-film?phrase=ocean)