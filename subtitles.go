package dgkdxf

import (
	"fmt"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"os"
	"strings"
)

type Subtitles struct {
	Begin     int    `json:"begin"`
	End       int    `json:"end"`
	Separator string `json:"separator"`
	Words     string `json:"words"`
}

func ConvertSubtitles2SrtFormat(subtitlesList []*Subtitles, srtFile string) error {
	if len(subtitlesList) == 0 {
		return dgerr.ARGUMENT_NOT_VALID
	}

	var srtBuilder strings.Builder
	for i, subtitles := range subtitlesList {
		srtBuilder.WriteString(fmt.Sprintf("%d\n", i+1))
		srtBuilder.WriteString(fmt.Sprintf("%s --> %s\n", formatMilliSecond2SubtitlesTime(subtitles.Begin), formatMilliSecond2SubtitlesTime(subtitles.End)))
		srtBuilder.WriteString(subtitles.Words)
		srtBuilder.WriteString("\n\n")
	}

	srtContent := srtBuilder.String()
	strings.TrimSuffix(srtContent, "\n")
	return os.WriteFile(srtFile, []byte(srtContent), os.ModePerm)
}

func formatMilliSecond2SubtitlesTime(milliSecond int) string {
	h := milliSecond / 1000 / 60 / 60
	m := milliSecond / 1000 / 60 % 60
	s := milliSecond / 1000 % 60
	ms := milliSecond % 1000

	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, ms)
}
