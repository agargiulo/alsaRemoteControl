package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	alsa "src.doom.fm/agargiulo/alsaRemoteControl"
)

func volumeStatusResponse(w http.ResponseWriter) {
	alsaVolume, err := alsa.GetVolume()
	if err != nil {
		panic(err)
	}
	isMuted, err := alsa.GetMuted()
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintf(w, "Volume is %d\nMuted: %t\n", alsaVolume, isMuted); err != nil {
		panic(err)
	}
}

func errBadRequestVolume(w http.ResponseWriter, err error) {
	message := "400 - [valid] POST /volume/(0-100)\n" + err.Error() + "\n"
	http.Error(w, string(message), http.StatusBadRequest)
}

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		http.Error(w, "GET or HEAD only", http.StatusMethodNotAllowed)
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
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		http.Error(w, "GET or HEAD only", http.StatusMethodNotAllowed)
		return
	}

	volumeStatusResponse(w)
}

func setVolume(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	volumeBase := path.Base(req.URL.Path)

	reqVolume, err := strconv.Atoi(volumeBase)
	if err != nil {
		errBadRequestVolume(w, err)
		return
	}
	err = alsa.SetVolume(reqVolume)
	if err != nil {
		errBadRequestVolume(w, err)
		return
	}
	volumeStatusResponse(w)
}

func toggle(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		http.Error(w, "GET or HEAD only", http.StatusMethodNotAllowed)
		return
	}
	err := alsa.Toggle()
	if err != nil {
		panic(err)
	}
	volumeStatusResponse(w)
}

func volUp(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		http.Error(w, "GET or HEAD only", http.StatusMethodNotAllowed)
		return
	}
	if err := alsa.IncreaseVolume(5); err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintln(w, "Volume went up by 5."); err != nil {
		panic(err)
	}
	volumeStatusResponse(w)
}

func volDown(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		http.Error(w, "GET or HEAD only", http.StatusMethodNotAllowed)
		return
	}

	if err := alsa.IncreaseVolume(-5); err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintln(w, "Volume went down by 5."); err != nil {
		panic(err)
	}
	volumeStatusResponse(w)
}

func mute(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		http.Error(w, "GET or HEAD only", http.StatusMethodNotAllowed)
		return
	}
	err := alsa.Mute()
	if err != nil {
		panic(err)
	}
	volumeStatusResponse(w)
}

func unmute(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		http.Error(w, "GET or HEAD only", http.StatusMethodNotAllowed)
		return
	}
	err := alsa.Unmute()
	if err != nil {
		panic(err)
	}
	volumeStatusResponse(w)
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
	if len(os.Args) >= 2 {
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
