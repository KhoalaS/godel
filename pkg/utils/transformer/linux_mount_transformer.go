package transformer

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/rs/zerolog/log"
)

var driveRegex = regexp.MustCompile(`^([A-Z]):/`)

func LinuxMountTransformer(job *types.DownloadJob) error {

	if runtime.GOOS != "linux" {
		log.Warn().Msg("Tried using linux_mount transformer on non Linux OS")
		return nil
	}

	m := driveRegex.FindStringSubmatch(job.Filename)
	if len(m) != 2 {
		log.Debug().Str("filename", job.Filename).Str("transformer", "linux_mount").Send()
		return fmt.Errorf("could not find drive name in filename %s", job.Filename)
	}

	match := m[0]
	driveChar := strings.ToLower(m[1])

	newFilename := strings.Replace(job.Filename, match, fmt.Sprintf("/mnt/%s/", driveChar), 1)
	job.Filename = newFilename

	return nil
}
