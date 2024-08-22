package dto

type Tune struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	RhythmFile string `json:"rhythm_file"`
	Owner      string `json:"owner"`
	CreatedAt  int64  `json:"created_at"`
}
