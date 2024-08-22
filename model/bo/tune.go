package bo

import "time"

type Tune struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	RhythmFile string    `json:"rhythm_file"`
	Owner      string    `json:"owner"`
	CreatedAt  time.Time `json:"created_at"`
}
