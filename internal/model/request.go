package model

import "task-api/internal/constant"

type RequestItem struct {
	Title    string  `binding:"required"`
	Price    float64 `binding:"gte=0"`
	Quantity uint `binding:"gte=0"`
	Owner string
}

type RequestFindItem struct {
	Statuses constant.ItemStatus `form:"status"`
	Title    string `form:"title"`


}

type RequestUpdateItem struct {
	Status constant.ItemStatus
}

type RequestUpdateIteminfo struct {
	Title    string  `json:"title"`
    Price    float64 `json:"price"`
    Quantity uint     `json:"quantity"`
}

type RequestLogin struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}
