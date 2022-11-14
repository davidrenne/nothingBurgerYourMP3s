import numpy
import struct
import os


class createEnv:

    def create_test_dir(self):
        os.mkdir('test_tmp')
        os.chdir('test_tmp')

    def create_audio_file(self):
        sampling_rate = 44100
        frequency = 440
        sample = 44100/440
        interval = numpy.arange(sample)