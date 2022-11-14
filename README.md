# cboMP3

Tool to re-encode recursively mp3 files to any bitrate.

## Install python dependencies with pip

```bash
pip install ffmpy
pip install mutagen
```

If on windows, the ffmpy seems to have binaries embedded for ffmpeg I believe.  If you are on mac do this `brew install ffmpeg` and on linux make sure the binary is in your path prior to running.

## Usage

By default, this will re-encode at 128 kbps because I cant hear the difference between that and higher fidelity. Only CD waves I can hear the difference

To start, clone the repo:

```bash
git clone https://github.com/davidrenne/nothingBurgerYourMP3s
cd nothingBurgerYourMP3s

go run main.go Y:\YourMP3s\ 384

# Or leave blank for 128

go run main.go Y:\YourMP3s\
```

## Warning

This tool will overwrite your files. Make sure you have backups prior to using

## Contributing

This program ideally should be a go install, but it has dependencies with python to pass to ffmpeg abstraction which calls either ffmpeg.exe or ffmpeg.  Ideally we can just use pure go for this someday.
