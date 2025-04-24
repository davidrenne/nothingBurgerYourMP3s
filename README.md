# nothingBurgerYourMP3s

Tool to re-encode recursively mp3 files to any bitrate. I used this to re-encode 2TB of mp3s. It "nothing burgered" 18,966 384kbps MP3 files and went from 2,379,714,605,056 bytes to 1,275,793,204,252 bytes. There are more files to process still but this is pretty solid right now.

<!-- TODO Remove python dependencies, make it all golang-->
<!-- TODO what it ffmpy and mutagen? -->
## Install python dependencies with pip 

```bash
pip install ffmpy
pip install mutagen
```

If on windows, the ffmpy seems to have binaries embedded for ffmpeg I believe. If you are on mac do this `brew install ffmpeg` and on linux make sure the binary is in your path prior to running.

## Usage

By default, this will re-encode at 128 kbps because I cant hear the difference between that and higher fidelity. Only CD waves I can hear the difference

To start, clone the repo:

```bash
git clone https://github.com/davidrenne/nothingBurgerYourMP3s
cd nothingBurgerYourMP3s

go run main.go Y:\YourMP3sOrFlacFiles\ 384

# Or leave blank for 128

go run main.go Y:\YourMP3sOrFlacFiles\
```

## Warning

This tool will overwrite your files. Make sure you have backups prior to using for safety reasons, but it should be fine. I have used it on thousands of files successfully.

## Contributing

This program ideally should be a go install, but it has dependencies with python to pass to ffmpeg abstraction which calls either ffmpeg.exe or ffmpeg. Ideally we can just use pure go for this someday. 