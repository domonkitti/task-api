package model

import "task-api/internal/constant"

type RequestItem struct {
	Title    string  `binding:"required"`
	Amount    float64 
	//`binding:"gte=5"`
	Quantity uint
}

type RequestFindItem struct {
	Statuses constant.ItemStatus `form:"status"`
}

type RequestUpdateItem struct {
	Status constant.ItemStatus
}

type RequestUpdateIteminfo struct {
	Title    string  `json:"title"`
    Amount    float64 `json:"Amount"`
    Quantity uint     `json:"quantity"`
}

type RequestLogin struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}
