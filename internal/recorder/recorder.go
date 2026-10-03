package recorder

import (
	"context"
	"fmt"
	"os/exec"
)

type Recorder struct {
	command *exec.Cmd

	identifier, endpoint, output string
}

func New(identifier, endpoint, output string) *Recorder {
	return &Recorder{identifier: identifier, endpoint: endpoint, output: output}
}

func (r *Recorder) Run(ctx context.Context) error {
	r.command = exec.CommandContext(
		ctx,
		"ffmpeg",

		"-use_wallclock_as_timestamps", "1",

		"-i", r.endpoint, "-an",

		"-codec:v", "libx265", "-tag:v", "hvc1",
		"-preset", "veryfast", "-tune", "zerolatency", "-crf", "32",

		"-g", "30", "-keyint_min", "30",

		"-vf", fmt.Sprintf("fps=5,drawtext=text='%%{localtime} %s':fontsize=32:fontcolor=white:x=32:y=32", r.identifier),

		"-f", "hls", "-hls_time", "6", "-hls_list_size", "6000", "-hls_flags", "append_list+delete_segments", "-hls_segment_type", "fmp4",

		r.output,
	)

	if err := r.command.Start(); err != nil {
		return err
	}

	return nil
}

func (r *Recorder) Stop() error {
	if err := r.command.Cancel(); err != nil {
		return err
	}

	return r.command.Process.Release()
}
