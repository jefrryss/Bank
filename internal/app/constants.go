package app

const (
	stateMain = iota
	stateImportSub
	stateExportSub
	stateAnalyticsSub
	stateBalanceSub

	stateManualEntityMenu
	stateManualActionMenu
	stateManualFormInput
	stateManualDeleteIDInput

	stateImportPathInput
	stateExportNameInput
	stateLogs
	stateShowMessage
	stateShowImportErrors
	stateAnalyticsDateFromInput
	stateAnalyticsDateToInput
	stateManualAccountIDInput
	stateManualNewBalanceInput
)

var mainMenuItems = []string{
	"Импорт",
	"Экспорт",
	"Пересчёт баланса",
	"Аналитика",
	"Ручной ввод",
	"Просмотр логов",
	"Выход",
}

var importSubItems = []string{
	"Импорт CSV",
	"Импорт JSON",
	"Просмотр ошибок при импорте",
	"Назад",
}

var exportSubItems = []string{
	"Экспорт CSV",
	"Экспорт JSON",
	"Назад",
}

var analyticsSubItems = []string{
	"Подсчет разницы доходов и расходов за период",
	"Группировка по категориям",
	"Назад",
}

var balanceSubItems = []string{
	"Автоматически",
	"Вручную",
	"Назад",
}

var manualEntityItems = []string{
	"Счета",
	"Категории",
	"Операции",
	"Назад",
}

var manualActionItems = []string{
	"Создать",
	"Редактировать",
	"Удалить",
	"Назад",
}
