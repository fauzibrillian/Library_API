package handler

type BookRequest struct {
	Tittle    string `json:"tittle" form:"tittle"`
	Publisher string `json:"publisher" form:"publisher"`
	Author    string `json:"author" form:"author"`
	Picture   string `json:"picture" form:"picture"`
}

type BookResponse struct {
	ID        uint   `json:"book_id"`
	Tittle    string `json:"tittle" form:"tittle"`
	Publisher string `json:"publisher" form:"publisher"`
	Author    string `json:"author" form:"author"`
	Picture   string `json:"picture" form:"picture"`
}

type BookPutRequest struct {
	ID        uint   `json:"book_id"`
	Tittle    string `json:"tittle" form:"tittle"`
	Publisher string `json:"publisher" form:"publisher"`
	Author    string `json:"author" form:"author"`
	Picture   string `json:"picture" form:"picture"`
}

type BookPutResponse struct {
	ID        uint   `json:"book_id"`
	Tittle    string `json:"tittle" form:"tittle"`
	Publisher string `json:"publisher" form:"publisher"`
	Author    string `json:"author" form:"author"`
	Picture   string `json:"picture" form:"picture"`
}
