package transcoding

import "os/exec"

func buildFFmpegCommand(input string, outDir string) *exec.Cmd {
	return exec.Command(
		"ffmpeg",
		"-i", input,
		"-map", "0:v:0",
		"-map", "0:a:0?",
		"-f", "dash",
		outDir+"/manifest.mpd",
	)
}
