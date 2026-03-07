package serializer

import "github.com/jhw66/myvideo_lab4/model"

type Vedio struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Info      string `json:"info"`
	CreatedAt int64  `json:"created_at"`
}

func BuildVedio(vedio *model.Vedio) Vedio {
	return Vedio{
		ID:        vedio.ID,
		Title:     vedio.Title,
		URL:       vedio.URL,
		Info:      vedio.Info,
		CreatedAt: vedio.CreatedAt.Unix(),
	}
}
