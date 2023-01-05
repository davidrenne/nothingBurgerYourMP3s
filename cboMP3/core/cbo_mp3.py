import glob
import ffmpy
import os
import sys
from mutagen.mp3 import MP3
import shutil
from sys import platform


class CBOMp3:
    def __init__(self, bitrate=128):
        if bitrate == None:
            bitrate = 128
        self.bitrate = bitrate

    @staticmethod
    def get_mp3_files(dir_path=""):
        """
        Scan selected directory for search all mp3 files
        dir_path: path for selected directory
        :return: list of mp3 files
        """
        return glob.glob("{}*.mp3".format(dir_path))

    def check_file_to_convert(self, mp3_files):
        """
        check bitrate of mp3 file and remove
        :return:
        """
        for mp3_file in mp3_files:
            audio = MP3(mp3_file)
            if audio.info.bitrate > (self.bitrate * 1000):
                print('File to convert: {}'.format(mp3_file))
            else:
                mp3_files.remove(mp3_file)
        return mp3_files

    def convert_files(self, mp3_files):
        """

        :param mp3_files:
        :return:
        """
        for mp3_file in mp3_files:
            if "wip-nothingburger-" in mp3_file:
                print(mp3_file + "is a previous nothingburger file deleting tmp file")
                os.remove(dst_mp3)
                exit(1)
            dst_mp3 = self.generate_tmp_file_location(mp3_file)
            cmd = ''
            if ".flac" in mp3_file:
                cmd = '-c:v copy'
            else:
                cmd = '-acodec libmp3lame'
            ff = ffmpy.FFmpeg(
                inputs={mp3_file: ''},
                outputs={dst_mp3: ' ' + cmd + ' -b:a ' +
                         str(self.bitrate) + 'k'}
            )

            try:
                ff.run()
                print('File {} enocoded.'.format(dst_mp3))
                shutil.copyfile(dst_mp3, mp3_file)
                if ".flac" in mp3_file:
                    shutil.move(mp3_file, mp3_file.replace(".flac", ".mp3"))
                print('File {} copied.'.format(mp3_file))
                os.remove(dst_mp3)
                print('tmp file {} deleted.'.format(dst_mp3))
            except ffmpy.FFRuntimeError:
                print('File {} already exists.'.format(dst_mp3))
            except Exception as e:
                print('Exception generic on ' + mp3_file, type(e), e)

    def generate_tmp_file_location(self, mp3_file):
        """

        :param mp3_file:
        :return:
        """
        folder = "/"
        if platform == "win32":
            folder = "\\"
        folders = mp3_file.split(folder)
        song = folders[len(folders) - 1]
        songPieces = song.split(".")
        if songPieces[len(songPieces) - 1] == "flac":
            songPieces[len(songPieces) - 1] = "mp3"
            song = ".".join([str(x) for x in songPieces])
        folders.pop()
        dst_file = folder.join([str(x) for x in folders])
        return dst_file + folder + "{}".format("wip-nothingburger-" + song)


if __name__ == "__main__":
    cbo = CBOMp3()
    if len(sys.argv) >= 2:
        bitrate = None
        if len(sys.argv) == 3:
            bitrate = sys.argv[2]
        cbo = CBOMp3(bitrate)
        mp3_files = [sys.argv[1]]
    else:
        # Hard code your recursive path to mp3 files
        mp3_files = cbo.get_mp3_files('J:\\myLargeMp3sIwantToSaveOnDiskSpace\\')
    cbo.convert_files(mp3_files)
