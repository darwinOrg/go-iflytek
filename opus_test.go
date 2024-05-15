package dgkdxf

import (
	"github.com/darwinOrg/go-common/utils"
	"os"
	"testing"
)

func TestExtractRawOpusData(t *testing.T) {
	data, _ := os.ReadFile("1.opus")
	rawData := ExtractRawOpusData(data)
	_ = utils.AppendToFile("1.raw.opus", rawData)
}
