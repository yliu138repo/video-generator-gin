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