package jsoner

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Jsoner struct {
	Data    map[string]any
	Mu      sync.Mutex
	Encoder *json.Encoder
}

func NewJson(file *os.File) *Jsoner {
	return &Jsoner{
		Data:    make(map[string]any),
		Encoder: json.NewEncoder(file),
	}
}

func (j *Jsoner) Write() {
	err := j.Encoder.Encode(j.Data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Не удалось записать данные в report.json")
		os.Exit(1)
	}
}

func (j *Jsoner) AddFileName(filename string) {
	j.Mu.Lock()
	defer j.Mu.Unlock()
	j.Data[filename] = make(map[string]any)
}

func (j *Jsoner) AddFileMetric(filename, metric string, value any) {
	j.Mu.Lock()
	defer j.Mu.Unlock()
	filedata, ok := j.Data[filename].(map[string]any)
	if !ok {
		fmt.Fprintln(os.Stderr, "gokode: File does not registered")
		return
	}
	filedata[metric] = value
}
