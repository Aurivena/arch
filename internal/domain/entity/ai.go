package entity

const (
	QwQRole = "user"
)

type AiSend struct {
	Model    string      `json:"model"`
	Messages []AiMessage `json:"messages"`
}

type AiMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type AiPlace struct {
	Title   string `json:"title"`
	Address string `json:"address"`
}
