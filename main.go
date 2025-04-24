package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/path"
	"github.com/DanielRenne/GoCore/core/utils"
)

type filesSync struct {
	sync.Mutex
	Items []string
}

func RecurseFiles(fileDir string) (files []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	var wg sync.WaitGroup
	var syncedItems filesSync
	path := fileDir

	if extensions.DoesFileExist(path) == false {
		return
	}

	err = filepath.Walk(path, func(path string, f os.FileInfo, errWalk error) (err error) {

		if errWalk != nil {
			err = errWalk
			return
		}

		if !f.IsDir() {
			wg.Add(1)
			syncedItems.Lock()
			syncedItems.Items = append(syncedItems.Items, path)
			syncedItems.Unlock()
			wg.Done()
		}

		return
	})
	wg.Wait()
	files = syncedItems.Items

	return
}

type processJob struct {
	Func func(string)
	File string
	Wg   *sync.WaitGroup
}

var (
	lockMp3sDone          sync.RWMutex
	jobs                  chan processJob
	mp3sDone              []string
	mp3sProcessedFileName string
)

func init() {
	numConcurrent := 10
	jobs = make(chan processJob)
	mp3sDone = make([]string, 0)
	mp3sProcessedFileName = "allMP3s.json"
	extensions.ReadFileAndParse(mp3sProcessedFileName, &mp3sDone)
	for i := 0; i < numConcurrent; i++ {
		go worker(i)
	}
}

func worker(idx int) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	for job := range jobs {
		job.Func(job.File)
		job.Wg.Done()
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please pass your MP3 directory to process")
	}
	var directoryToIterate string
	var bitRate string
	potentialPath := os.Args[1]
	if len(os.Args) == 3 {
		bitRate = os.Args[2]
	} else {
		bitRate = "128"
	}
	lastByte := potentialPath[len(potentialPath)-1:]
	if lastByte != "\\" && path.IsWindows {
		directoryToIterate = potentialPath + "\\"
	} else if lastByte != "/" {
		directoryToIterate = potentialPath + "/"
	}

	if path.IsWindows && strings.Index(directoryToIterate, "\\\\") != -1 {
		log.Fatal("Please only escape your directory path once with \\")
	}

	if extensions.DoesFileExist(directoryToIterate) == false {
		log.Fatal("Path does not exist or is invalid")
	}

	var wg sync.WaitGroup
	startEntireProcess := time.Now()
	var processJobs []processJob
	files, _ := RecurseFiles(directoryToIterate)
	for _, fileToWorkOn := range files {
		pieces := strings.Split(fileToWorkOn, ".")
		ext := strings.ToUpper(pieces[len(pieces)-1:][0])
		if ext == "MP3" || ext == "FLAC" {
			processJobs = append(processJobs, processJob{
				Wg:   &wg,
				File: fileToWorkOn,
				Func: func(fileWork string) {
					start := time.Now()
					lockMp3sDone.RLock()
					if utils.InArray(fileWork, mp3sDone) {
						log.Println("Skipping: " + fileWork)
						lockMp3sDone.RUnlock()
						return
					}
					lockMp3sDone.RUnlock()

					intBitRate, err := strconv.Atoi(bitRate)
					if err != nil {
						log.Println(fileWork, " Errored\n\n\n", err.Error())
						return
					}
					err = convertMusicFile(fileWork, intBitRate)
					if err != nil {
						log.Println(fileWork, " Errored\n\n\n", err.Error())
						return
					}

					lockMp3sDone.Lock()
					mp3sDone = append(mp3sDone, fileWork)
					data, err := json.MarshalIndent(mp3sDone, "", "    ")
					if err == nil {
						err = extensions.Write(string(data), mp3sProcessedFileName)
					}
					lockMp3sDone.Unlock()

					log.Println(fileWork, " Done! \n\nTook "+logger.TimeTrack(start, "func time")+"\n\n")

					/*
						stdOut, stdErr, err := cmdExec.Run("python3", "cboMP3/core/cbo_mp3.py", fileWork, bitRate)
						if err != nil {
							log.Println(fileWork, " Errored\n\n\n", stdOut, stdErr, err.Error())
							return
						}

						lockMp3sDone.Lock()
						mp3sDone = append(mp3sDone, fileWork)
						data, err := json.MarshalIndent(mp3sDone, "", "    ")
						if err == nil {
							err = extensions.Write(string(data), mp3sProcessedFileName)
						}
						lockMp3sDone.Unlock()

						log.Println(fileWork, " Done! \n\nTook "+logger.TimeTrack(start, "func time")+"\n\n", stdOut, stdErr)
					*/
				},
			})
		}
	}

	wg.Add(len(processJobs))
	go func() {
		for _, job := range processJobs {
			j := job
			jobs <- j
		}
	}()

	logger.Log("Waiting on all " + extensions.IntToString(len(processJobs)) + " MP3 ffmpeg go routines to finish...")
	wg.Wait()
	log.Println(logger.TimeTrack(startEntireProcess, "Completed in"))
}

// to remove python

func findFFmpegExecutable() (string, error) {
	// TODO - do we need to find ffmpeg.exe in the PATH?
	// Try to find ffmpeg.exe in the PATH
	// ffmpegPath, err := exec.LookPath("ffmpeg.exe")
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err == nil {
		return ffmpegPath, nil
	}

	// If not found in PATH, try the program files directory on Windows
	if thisOS := strings.ToLower(runtime.GOOS); thisOS == "windows" {
		programFiles := os.Getenv("ProgramFiles")
		ffmpegPath = filepath.Join(programFiles, "ffmpeg", "bin", "ffmpeg.exe")
		_, err = os.Stat(ffmpegPath)
		if err == nil {
			return ffmpegPath, nil
		}
	}

	return "", fmt.Errorf("ffmpeg not found")
}

func convertMusicFile(filename string, bitRate int) error {
	if bitRate <= 0 {
		bitRate = 128
	}

	ffmpegPath, err := findFFmpegExecutable()
	if err != nil {
		return fmt.Errorf("FFmpeg not found: %v", err)
	}

	// Check if the file exists
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		return fmt.Errorf("File does not exist: %s", filename)
	}

	// Get the file extension
	ext := strings.ToLower(filepath.Ext(filename))

	// Check if it's an MP3 or FLAC file
	isFLAC := ext == ".flac"
	isMP3 := ext == ".mp3"

	if !isFLAC && !isMP3 {
		return fmt.Errorf("Unsupported file format: %s", ext)
	}

	// Check the bitrate of the input file
	if isMP3 {
		cmd := exec.Command(ffmpegPath, "-i", filename)
		output, _ := cmd.CombinedOutput()
		// TODO ignore this error or use ffprobe instead
		// if err != nil {
		// 	log.Println(string(output))
		// 	return fmt.Errorf("Error checking bitrate: %v", err)
		// }
		bitrateStr := string(output)
		if strings.Contains(bitrateStr, "Audio: mp3") {
			currentBitrate := parseBitrate(bitrateStr)
			if currentBitrate <= bitRate {
				return fmt.Errorf("Bitrate is already lower than the desired output: %d kbps", currentBitrate)
			}
		}
	}

	// Generate the output filename
	outputFilename := "wip-nothingburger-" + filepath.Base(filename)
	// outputFile := filepath.Join(filepath.Dir(filename), outputFilename)

	// Determine the codec and options based on the input file format
	codec := "-c:a libmp3lame"
	if isFLAC {
		codec = "-c:a flac"
	}

	// Run FFmpeg command
	log.Println(ffmpegPath, "-i", filename, codec, "-b:a", fmt.Sprintf("%dk", bitRate), "-y", outputFilename)
	cmd := exec.Command(ffmpegPath, "-i", filename, "-b:a", fmt.Sprintf("%dk", bitRate), "-y", outputFilename)
	// ffmpeg -i input.mp3 -codec:a libmp3lame -b:a 128k output.mp3

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Println(outputFilename)
		return fmt.Errorf("Error converting file: %v", err)
	}

	// Replace the original file with the converted file
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("Error removing original file: %v", err)
	}

	err = os.Rename(outputFilename, filename)
	if err != nil {
		return fmt.Errorf("Error renaming the converted file: %v", err)
	}

	fmt.Printf("File %s encoded.\n", filename)
	return nil
}

func parseBitrate(output string) int {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Stream #0:0") {
			parts := strings.Fields(line)
			for i := 0; i < len(parts); i++ {
				if parts[i] == "kb/s" {
					bitRate, err := strconv.Atoi(parts[i-1])
					if err != nil {
						fmt.Errorf("Error while parsing bitrate: %v", err)
					}
					return bitRate
				}
			}
		}
	}
	return 0
}
