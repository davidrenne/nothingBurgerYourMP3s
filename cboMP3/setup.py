import os
from setuptools import setup


def read(fname):
    return open(os.path.join(os.path.dirname(__file__), fname)).read()


setup(
    name="cboMP3",
    version="1.0.0d",
    author="Madaoo",
    author_email="kielkowskimarcin@prokonto.pl",
    description=("Change bitrate in mp3 file"),
    license="MIT",
    keywords="mp3 bitratw converter",
    packages=['mutagen', 'ffmpy'],
    url='https://github.com/MadaooQuake/cboMP3',
    classifiers=[
        "Development:w"
        "Status :: 1.0.0dev1",
        "Topic ::  app",
        "License :: OSI Approved :: BSD License",
    ], install_requires=['ffmpy', 'mutagen', 'numpy']
)
