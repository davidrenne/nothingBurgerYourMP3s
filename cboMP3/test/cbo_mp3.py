import unittest

from core.cbo_mp3 import CBOMp3
from test.create_test_env import createEnv


class CboMp3Test(unittest.TestCase):

    def setUp(self):
        pass

    def test_1_get_mp3_files(self):
        self.assertEqual(CBOMp3().get_mp3_files("C:"), [])

    def test_2_get_mp3_files(self):
        self.assertEqual(CBOMp3().get_mp3_files("aaa"), [])

    def test_3_get_mp3_files(self):
        self.assertEqual(CBOMp3().get_mp3_files(""), [])

    def test_4_get_mp3_files(self):
        self.assertEqual(CBOMp3().get_mp3_files(), [])

    def test_5_get_mp3_files(self):
        createEnv().create_test_dir()
        self.assertEquals(1, 2)

    def test_1_check_file_to_convert(self):
        self.assertEquals(1, 2)

    def test_1_convert_files(self):
        self.assertEquals(1, 2)

    def test_1_create_dst_location(self):
        self.assertEquals(1, 2)

    def test_1_create_tmp_dir(self):
        self.assertEquals(1, 2)




if __name__ == '__main__':
    unittest.main()
