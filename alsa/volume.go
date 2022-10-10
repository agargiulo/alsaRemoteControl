package alsa

// Most of this was copied in part from here:
// github.com/itchyny/volume-go
// and modified to better fit my needs

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	AMIXER = "amixer"
	// TODO: Make this configurable
	SCONTROLLER = "Digital"
)

var (
	amixerBaseCmd = []string{AMIXER, "-M"}
	volumePattern = regexp.MustCompile(`\d+%`)
)

func parseVolume(out string) (int, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		s := strings.TrimLeft(line, " \t")
		if strings.Contains(s, "Playback") && strings.Contains(s, "%") {
			volumeStr := volumePattern.FindString(s)
			return strconv.Atoi(volumeStr[:len(volumeStr)-1])
		}
	}
	return 0, errors.New("no volume found")
}

func parseMuted(out string) (bool, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		s := strings.TrimLeft(line, " \t")
		if strings.Contains(s, "Playback") && strings.Contains(s, "%") {
			if strings.Contains(s, "[off]") {
				return true, nil
			} else if strings.Contains(s, "[on]") {
				return false, nil
			}
		}
	}
	return false, errors.New("no mute information found")
}

func getVolumeCmd() []string {
	return append(amixerBaseCmd, []string{"get", SCONTROLLER}...)
}

func setVolumeCmd(volume int) []string {
	return append(amixerBaseCmd, []string{"set", SCONTROLLER, strconv.Itoa(volume) + "%"}...)
}

func increaseVolumeCmd(diff int) []string {
	var sign string
	if diff >= 0 {
		sign = "+"
	} else {
		// Thank you for having 5%+ and 5%- for amixer ... just thanks
		diff = -diff
		sign = "-"
	}
	return append(amixerBaseCmd, []string{"set", SCONTROLLER, strconv.Itoa(diff) + "%" + sign}...)
}

func getMutedCmd() []string {
	return append(amixerBaseCmd, []string{"get", SCONTROLLER}...)
}

func muteCmd() []string {
	return append(amixerBaseCmd, []string{"set", SCONTROLLER, "mute"}...)
}

func unmuteCmd() []string {
	return append(amixerBaseCmd, []string{"set", SCONTROLLER, "unmute"}...)
}

func toggleCmd() []string {
	return append(amixerBaseCmd, []string{"set", SCONTROLLER, "toggle"}...)
}

func cmdEnv() []string {
	return []string{"LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8"}
}

func execCmd(cmdArgs []string) ([]byte, error) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = append(os.Environ(), cmdEnv()...)
	out, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf(`failed to execute "%v" (%+v)`, strings.Join(cmdArgs, " "), err)
	}
	return out, err
}

// GetVolume returns the current volume (0 to 100).
func GetVolume() (int, error) {
	out, err := execCmd(getVolumeCmd())
	if err != nil {
		return 0, err
	}
	return parseVolume(string(out))
}

// SetVolume sets the sound volume to the specified value (0 to 100).
func SetVolume(volume int) error {
	if volume < 0 || 100 < volume {
		return errors.New("out of valid (0-100) volume range")
	}
	_, err := execCmd(setVolumeCmd(volume))
	return err
}

// IncreaseVolume increases (or decreases) the audio volume by the specified value.
func IncreaseVolume(diff int) error {
	_, err := execCmd(increaseVolumeCmd(diff))
	return err
}

// GetMuted returns the current muted status
func GetMuted() (bool, error) {
	out, err := execCmd(getMutedCmd())
	if err != nil {
		return false, err
	}
	return parseMuted(string(out))
}

func Mute() error {
	_, err := execCmd(muteCmd())
	return err
}
func Unmute() error {
	_, err := execCmd(unmuteCmd())
	return err
}

func Toggle() error {
	_, err := execCmd(toggleCmd())
	return err
}
