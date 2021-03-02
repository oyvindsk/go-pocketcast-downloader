package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/oyvindsk/go-id3v2" // forked from github.com/bogem/id3v2
	// github.com/mikkyang/id3-go Writes invalid files according to kid3-cli and DeaDBeeF
)

func downloadAllMissing(starred *resultStarred, downloadedFilesPath, successFilepath string) error {

	// Ensure the output dir exists and is a dir
	d, err := os.Stat(downloadedFilesPath)

	// exits, but not a dir
	if err == nil && !d.IsDir() {
		return fmt.Errorf("downloadAllMissing: output dir exists, but is not a directory: %q", downloadedFilesPath)
	}

	if err != nil {

		// error, other than doesNotExist
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("downloadAllMissing: stating output dir: %q, err: %s", downloadedFilesPath, err)
		}

		// does not exists - Create it
		log.Printf("Creating output dir: %q", downloadedFilesPath)
		err = os.Mkdir(downloadedFilesPath, 0770)
		if err != nil {
			return fmt.Errorf("downloadAllMissing: creating output dir: %q, err: %s", downloadedFilesPath, err)
		}
	}

	// Load the file containing uuids we already have
	gotAlready, err := loadSuccessfullyDownloaded(successFilepath)
	if err != nil {
		return fmt.Errorf("downloadAllMissing: loadSuccessfulyDownloaded: %s", err)
	}

	var successCnt = len(gotAlready) // int // , failedCnt int

	for _, e := range starred.Episodes {

		if _, found := gotAlready[e.UUID]; found {
			continue
		}

		fpath := strings.Replace(fmt.Sprintf("%s -- %s.mp3", e.PodcastTitle, e.Title), string(os.PathSeparator), " - ", 100)
		fpath = filepath.Join(downloadedFilesPath, fpath)
		fpath = filepath.Clean(fpath)

		log.Printf(`Getting:
			%q -- %q (%s)
			to file: %q`,
			e.PodcastTitle, e.Title, e.UUID, fpath,
		)

		err := DownloadFile(fpath, e.URL)
		if err != nil {
			return fmt.Errorf("downloadAllMissing: %s", err)
			// TODO : don't give up on all if there's a isolated error ?
		}

		// Tag it!
		err = tagFile(fpath, e.PodcastTitle, e.Title, true)
		if err != nil {
			for {
				e := errors.Unwrap(err)
				if e == nil {
					break
				}
				log.Printf("Unwraped err: %s", e)
			}

			return fmt.Errorf("downloadAllMissing: %s", err)
			// TODO : don't give up on all if there's a isolated error ?
		}

		successCnt++

		log.Printf(`********************* Episode completed! Success count: %d of %d *********************`, successCnt, len(starred.Episodes))

		// Success! Store the uuid so we don't fetch it again
		err = appendSuccessfullyDownloaded(successFilepath, e.UUID)
		if err != nil {
			log.Printf("downloadAllMissing: episode downloaded, but we cold not remember that we got it successfully. Err: %s", err)
		}
	}

	return nil
}

func loadSuccessfullyDownloaded(filepath string) (map[string]struct{}, error) {

	res := make(map[string]struct{}) // return empty if we can't find anything

	file, err := os.Open(filepath)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("loadSuccessfullyDownloaded: not found")
		return res, nil
	}
	if err != nil {
		return res, err
	}

	log.Printf("loadSuccessfullyDownloaded: found")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res[scanner.Text()] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func appendSuccessfullyDownloaded(filepath, uuid string) error {

	// If the file doesn't exist, create it, or append to the file
	// from the os docs - untested on windows and mac
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("appendSuccessfullyDownloaded: %s", err)
	}

	if _, err := f.WriteString(uuid + "\n"); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Printf("appendSuccessfullyDownloaded: writing: %s", err)
	}
	if err := f.Close(); err != nil {
		log.Printf("appendSuccessfullyDownloaded: closing: %s", err)
	}

	return nil
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("DownloadFile: %s", err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("DownloadFile: %s", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("DownloadFile: %s", err)
	}

	return nil
}

func tagFile(filepath, artist, title string, overwriteExisting bool) error {

	tag, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("tagFile: open: %s", err)
	}

	log.Printf("tagFile: default encoding: %s", tag.DefaultEncoding())

	log.Printf(`tagFile:
		was: Artist: %q, Album: %q, Title: %q`,
		tag.Artist(), tag.Album(), tag.Title(),
	)

	if overwriteExisting || tag.Artist() == "" {
		log.Printf("\tsetting Artist to %q", artist)
		tag.SetArtist(artist)
	}
	if overwriteExisting || tag.Title() == "" {
		log.Printf("\tsetting Title to %q", title)
		tag.SetTitle(title)
	}

	if overwriteExisting || tag.Genre() == "" {
		log.Printf("\tsetting Genre to %q", "Podcast")
		tag.SetGenre("Podcast")
	}

	err = tag.Save()
	if err != nil {
		return fmt.Errorf("tagFile: Save: %s", err)
	}

	return nil
}
