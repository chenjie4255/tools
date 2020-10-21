package avatarsample

import (
	"encoding/json"
	"github.com/chenjie4255/tools/slice"
	"os"
)

var samples []string

// GetSamples获取头像样本数
func GetSamples(count int) []string {
	idxs := slice.RandomIndexes(len(samples), count)
	ret := []string{}
	for _, idx := range idxs {
		ret = append(ret, samples[idx])
	}

	return ret
}

func LoadSample(filePath string) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	fileContent := struct {
		Avatars []string `json:"avatars"`
	}{}

	if err := json.NewDecoder(f).Decode(&fileContent); err != nil {
		return 0, err
	}

	samples = fileContent.Avatars
	return len(samples), nil
}
