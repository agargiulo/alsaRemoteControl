package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func volumeStatusResponse(w http.ResponseWriter) {
	alsaVolume, err := GetVolume()
	if err != nil {
		panic(err)
	}
	isMuted, err := GetMuted()
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintf(w, "Volume is %d\nMuted: %t", alsaVolume, isMuted); err != nil {
		panic(err)
	}
}

func errBadRequestVolume(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	message := []byte("400 - [valid] POST /volume/(0-100)\n")
	if err != nil {
		message = append(message, err.Error()+"\n"...)
	}
	if _, err := w.Write(message); err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}

	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	if _, err := fmt.Fprint(w, "Welcome!\n"); err != nil {
		panic(err)
	}
}

func volume(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}

	volumeStatusResponse(w)
}

func setVolume(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.NotFound(w, req)
		return
	}
	volumeBase := path.Base(req.URL.Path)

	reqVolume, err := strconv.Atoi(volumeBase)
	if err != nil {
		errBadRequestVolume(w, err)
		return
	}
	err = SetVolume(reqVolume)
	if err != nil {
		errBadRequestVolume(w, err)
		return
	}
	volumeStatusResponse(w)
}

func toggle(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}
	err := Toggle()
	if err != nil {
		panic(err)
	}
	muted, err := GetMuted()
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintf(w, "Volume is now muted: %t\n", muted); err != nil {
		panic(err)
	}
}

func volUp(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}
	if err := IncreaseVolume(5); err != nil {
		panic(err)
	}
	alsaVolume, err := GetVolume()
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintf(w, "Volume went up by 5. Current volume: %d\n", alsaVolume); err != nil {
		panic(err)
	}
}

func volDown(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}

	if err := IncreaseVolume(-5); err != nil {
		panic(err)
	}
	alsaVolume, err := GetVolume()
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintf(w, "Volume went down by 5. Current volume: %d\n", alsaVolume); err != nil {
		panic(err)
	}
}

func mute(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}
	err := Mute()
	if err != nil {
		panic(err)
	}
	muted, err := GetMuted()
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintf(w, "Audio is now muted: %t\n", muted); err != nil {
		panic(err)
	}
}

func unmute(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}
	err := Unmute()
	if err != nil {
		panic(err)
	}
	muted, err := GetMuted()
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintf(w, "Audio is now unmuted: %t\n", !muted); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/volume", volume)
	http.HandleFunc("/volume/", setVolume)
	http.HandleFunc("/mute", mute)
	http.HandleFunc("/unmute", unmute)
	http.HandleFunc("/toggle", toggle)
	http.HandleFunc("/up", volUp)
	http.HandleFunc("/down", volDown)

	var port string
	if len(os.Args) >= 1 {
		port = os.Args[1]
	} else {
		port = "12345"
	}
	certFile := os.Getenv("ALSA_REMOTE_SSL_CERT")
	certKey := os.Getenv("ALSA_REMOTE_SSL_KEY")
	if len(certFile) > 1 && len(certKey) > 1 {
		log.Fatal(http.ListenAndServeTLS(":"+port, certFile, certKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}
