package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/cmdExec"
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
					stdOut, stdErr, err := cmdExec.Run("python", "cboMP3/core/cbo_mp3.py", fileWork, bitRate)
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
