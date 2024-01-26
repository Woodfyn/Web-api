package domain

type Game struct {
	Id         int    `json:"id"`
	Title      string `json:"title" binding:"required"`
	Genre      string `json:"genre"`
	Evaluation int    `json:"evaluation"`
}

type UpdateItemInput struct {
	Id         *int    `json:"id"`
	Title      *string `json:"title" binding:"required"`
	Genre      *string `json:"genre"`
	Evaluation *int    `json:"evaluation"`
	Done       *bool   `json:"done"`
}
