package handler

type BookRequest struct {
	Tittle    string `json:"tittle" form:"tittle"`
	Publisher string `json:"publisher" form:"publisher"`
	Author    string `json:"author" form:"author"`
	Picture   string `json:"picture" form:"picture"`
	Category  string `json:"category" form:"category"`
	Stock     int    `json:"stock" form:"stock"`
}

type BookResponse struct {
	ID        uint   `json:"book_id"`
	Tittle    string `json:"tittle" form:"tittle"`
	Publisher string `json:"publisher" form:"publisher"`
	Author    string `json:"author" form:"author"`
	Picture   string `json:"picture" form:"picture"`
	Category  string `json:"category" form:"category"`
	Stock     int    `json:"stock" form:"stock"`
}

type BookPutRequest struct {
	ID        uint   `json:"book_id"`
	Tittle    string `json:"tittle" form:"tittle"`
	Publisher string `json:"publisher" form:"publisher"`
	Author    string `json:"author" form:"author"`
	Picture   string `json:"picture" form:"picture"`
	Category  string `json:"category" form:"category"`
}

type BookPutResponse struct {
	ID        uint   `json:"book_id"`
	Tittle    string `json:"tittle" form:"tittle"`
	Publisher string `json:"publisher" form:"publisher"`
	Author    string `json:"author" form:"author"`
	Picture   string `json:"picture" form:"picture"`
	Category  string `json:"category" form:"category"`
}

type BookDetailRequest struct {
	IdBook uint `json:"id_book" form:"id_book"`
	IdRack uint `json:"id_rack" form:"id_rack"`
}

type BookDetailResponse struct {
	IdBook     uint `json:"id_book" form:"id_book"`
	IdRack     uint `json:"id_rack" form:"id_rack"`
}
