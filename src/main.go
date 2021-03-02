package main

import (
    "log"
    "os"
)

const (
    downloadedFilesPath            = "./downloaded"
    authTokenFilepath              = "./authtoken.jwt"
    successfullyDownloadedFilepath = "./downloaded-successfully.txt"
)

func main() {

    // log.Fatalln(tagFile("./foo.mp3", "ART", "TTTTILE", false))

    // Login - look for token on disk
    found, authToken, err := loadAuthToken(authTokenFilepath)
    if err != nil {
        log.Fatalln(err)
    }

    // Login - if none, get one from the webserver using a user supplied username/pass
    if !found {

        // Read usernamer and password
        username, password := readUsernamePassword()
        if password == "" || username == "" {
            log.Fatalf("%s need the Pocketcast web username and passord in either the first and second command-line argument, or the environment variables: PCUSERNAME and PCPASSWORD", os.Args[0])
        }

        // Get jwt token using user/pass
        authToken, err = login(username, password)
        if err != nil {
            log.Fatalln(err)
        }

        // Store it, if successful
        err = storeAuthToken(authTokenFilepath, authToken)
        if err != nil {
            log.Fatalln(err)
        }
    }

    // Get starred podcasts (metadata) from Pocketcast
    starred, err := getStarred(authToken)
    if err != nil {
        log.Fatalln(err)
    }

    log.Printf("********************* Found %d starred episodes *********************", len(starred.Episodes))

    // log.Printf("got already:\n%+v", gotAlready)

    err = downloadAllMissing(starred, downloadedFilesPath, successfullyDownloadedFilepath)
    if err != nil {
        log.Println("********************* Downloading stopped, with error *********************")
        log.Fatalf("Some files downloaded, probably, but not all. Err: %s", err)
    }

    log.Println("********************* Success! Download completed! *********************")

}

func readUsernamePassword() (string, string) {
    var username, password string
    username = os.Getenv("PCUSERNAME")
    password = os.Getenv("PCPASSWORD")
    if (password == "" || username == "") && len(os.Args) == 3 {
        // try command line arguments instead
        username = os.Args[1]
        password = os.Args[2]
    }
    return username, password
}
