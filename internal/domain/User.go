package domain

type User struct {
	ID           string `bson:"_id,omitempty"     json:"id,omitempty"` // строковый ID
	FullName     string `bson:"fullName"          json:"fullName"`     // ФИО
	Password     string `bson:"password"          json:"password"`     // пароль (пока не хэшируется)
	Role         string `bson:"role"              json:"role"`         // "admin", "librarian", "reader"
	Phone        string `bson:"phone"             json:"phone"`        // телефон
	RegisteredAt string `bson:"registeredAt"      json:"registeredAt"` // дата регистрации (ISO string)
	IsActive     bool   `bson:"isActive"          json:"isActive"`     // активен или заблокирован
}

type UserFilter struct {
	FullNameContains string `json:"fullName"`   // фильтр по части ФИО
	Phone            string `json:"phone"`      // фильтр по телефону
	Role             string `json:"role"`       // фильтр по роли
	OnlyActive       *bool  `json:"onlyActive"` // null — все, true/false — по флагу активности
}
