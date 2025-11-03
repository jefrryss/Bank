package handmaker

type Field struct {
	Key         string
	Label       string
	Placeholder string
	Optional    bool
}

// Порядок и подсказки для форм
func AccountCreateFields() []Field {
	return []Field{
		{Key: "id", Label: "ID счёта", Placeholder: "например: 1"},
		{Key: "name", Label: "Название счёта", Placeholder: "например: Основной счёт"},
		{Key: "balance", Label: "Начальный баланс", Placeholder: "например: 120000.50"},
	}
}

func AccountUpdateFields() []Field {
	return []Field{
		{Key: "id", Label: "ID счёта", Placeholder: "существующий ID"},
		{Key: "name", Label: "Новое имя", Placeholder: "новое название"},
		{Key: "balance", Label: "Новый баланс", Placeholder: "например: 150000"},
	}
}

func CategoryCreateFields() []Field {
	return []Field{
		{Key: "id", Label: "ID категории", Placeholder: "например: 2"},
		{Key: "name", Label: "Название категории", Placeholder: "например: Питание"},
		{Key: "type", Label: "Тип (доход/расход)", Placeholder: "расход"},
	}
}

func CategoryUpdateFields() []Field {
	return []Field{
		{Key: "id", Label: "ID категории", Placeholder: "существующий ID"},
		{Key: "name", Label: "Новое название", Placeholder: "например: Еда"},
		{Key: "type", Label: "Тип (доход/расход)", Placeholder: "расход"},
	}
}

func OperationCreateFields() []Field {
	return []Field{
		{Key: "id", Label: "ID операции", Placeholder: "например: 10"},
		{Key: "type", Label: "Тип (доход/расход)", Placeholder: "расход"},
		{Key: "account_id", Label: "ID счёта", Placeholder: "существующий ID счёта"},
		{Key: "category_id", Label: "ID категории", Placeholder: "существующий ID категории"},
		{Key: "amount", Label: "Сумма", Placeholder: "например: 1999.99"},
		{Key: "date", Label: "Дата (YYYY-MM-DD)", Placeholder: "2024-01-31"},
		{Key: "description", Label: "Описание", Placeholder: "необязательно", Optional: true},
	}
}

func OperationUpdateFields() []Field {
	return []Field{
		{Key: "id", Label: "ID операции", Placeholder: "существующий ID"},
		{Key: "type", Label: "Новый тип (доход/расход)", Placeholder: "расход"},
		{Key: "account_id", Label: "Новый ID счёта", Placeholder: "существующий ID счёта"},
		{Key: "category_id", Label: "Новый ID категории", Placeholder: "существующий ID категории"},
		{Key: "amount", Label: "Новая сумма", Placeholder: "например: 2500"},
		{Key: "date", Label: "Новая дата (YYYY-MM-DD)", Placeholder: "2024-02-05"},
		{Key: "description", Label: "Новое описание", Placeholder: "необязательно", Optional: true},
	}
}
