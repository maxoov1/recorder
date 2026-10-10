package recorder

import (
	"context"
	"fmt"
	"os/exec"
)

type Recorder struct {
	command *exec.Cmd

	identifier string
	endpoint   string

	rollingStreamOutput string
	thumbnailOutput     string
}

func New(identifier, endpoint, rollingStreamOutput, thumbnailOutput string) *Recorder {
	return &Recorder{
		identifier: identifier,
		endpoint:   endpoint,

		rollingStreamOutput: rollingStreamOutput,
		thumbnailOutput:     thumbnailOutput,
	}
}

func (r *Recorder) Run(ctx context.Context) error {
	r.command = exec.CommandContext(
		ctx,
		"ffmpeg", "-y",
		"-use_wallclock_as_timestamps", "1", "-i", r.endpoint, "-an",
		"-codec:v", "libx265", "-tag:v", "hvc1",
		"-preset", "veryfast", "-tune", "zerolatency", "-crf", "32",
		"-g", "30", "-keyint_min", "30",
		"-filter_complex", fmt.Sprintf(`
			[0:v] drawtext=text='%%{localtime} %s':fontsize=32:fontcolor=white:x=32:y=32, split=2 [s][t];
			  [s] fps=5 [fs]; [t] fps=1, scale=640:-2 [ft]
			`, r.identifier),
		"-map", "[fs]", r.rollingStreamOutput,
		"-map", "[ft]", "-update", "1", "-q:v", "31", r.thumbnailOutput,
	)
	return r.command.Start()
}

func (r *Recorder) Stop() error {
	if err := r.command.Cancel(); err != nil {
		return err
	}

	_ = r.command.Wait()
	return nil
}
