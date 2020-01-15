package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
)

func volume(w http.ResponseWriter, req *http.Request) {
	alsaVolume, err := GetVolume()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Volume is %d\n", alsaVolume)
}

func setVolume(w http.ResponseWriter, req *http.Request) {
	volumeBase := path.Base(req.URL.Path)
	if volumeBase == "volume" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - choose either /volume or /volume/(0-100)\n"))
		return
	}
	reqVolume, err := strconv.Atoi(volumeBase)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Trying to set volume to: %d\n", reqVolume)
}

func toggle(w http.ResponseWriter, req *http.Request) {
	err := Toggle()
	if err != nil {
		panic(err)
	}
	muted, err := GetMuted()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Volume is now muted: %t\n", muted)
}

func volUp(w http.ResponseWriter, req *http.Request) {
	if err := IncreaseVolume(5); err != nil {
		panic(err)
	}
	alsaVolume, err := GetVolume()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Volume went up by 5. Current volume [%d]\n", alsaVolume)
}

func volDown(w http.ResponseWriter, req *http.Request) {
	if err := IncreaseVolume(-5); err != nil {
		panic(err)
	}
	alsaVolume, err := GetVolume()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Volume went down by 5. Current volume [%d]\n", alsaVolume)
}

func index(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/volume", volume)
	http.HandleFunc("/volume/", setVolume)
	http.HandleFunc("/toggle", toggle)
	http.HandleFunc("/up", volUp)
	http.HandleFunc("/down", volDown)

	log.Fatal(http.ListenAndServe(":12345", nil))
}
