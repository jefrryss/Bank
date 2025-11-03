<h1 align="center">TASK2 — Финансовый TUI-менеджер на Go</h1>
<p align="center">
Учет счетов, категорий и операций. Импорт/экспорт (CSV/JSON), аналитика, логи, DI (Uber Dig), TUI (Bubble Tea).
</p>

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.21%2B-00ADD8">
  <img alt="UI" src="https://img.shields.io/badge/TUI-bubbletea-6A4C93">
  <img alt="DI" src="https://img.shields.io/badge/DI-uber%2Fdig-1E90FF">
</p>

---

## Структура проекта

```text
task2/
├─ docs/                         # Хранит Readme   
├─ exampleFiles/                 # Файлы для импорта (примеры файлов)
├─ cmd/
│  └─ exportFiles/               # Сюда сохраняются файлы для экспорта
│  └─ main.go                    # Точка входа, DI-контейнер (uber/dig)
├─ domain/
│  ├─ entities/                  # Банковский счёт, Категория, Операция
│  ├─ factory/                   # Валидация/создание сущностей (BankAccount, Category, Operation)
│  ├─ ports/                     # Интерфейсы (DataImporter, ExporterVisitor)
│  ├─ repository/                # Интерфэйсы для репоиториев
│  └─ errordata/                 # ErrorRecord для записи ошибок
├─ internal/
│  ├─ app/                       # TUI-приложение (Bubble Tea) — состояния, экраны, навигация
│  ├─ bankmanager/               # Единая точка сохранения/поиска/удаления сущностей
│  ├─ balancemanager/            # Пересчет баланса (авто/вручную)
│  ├─ importer/                  # CSV/JSON импортеры + ImportFacade + лог декоратор
│  ├─ repository/                # Реализации интерфейсов репозиториев
│  ├─ export/                    # CSV/JSON экспорт + ExportFacade + визиторы + декоратор для логов + адаптеры
│  ├─ handmaker/                 # Ручной ввод c валидацией через фабрики
│  ├─ logstorage/                # Хранилище логов
│  └─ analitic/                  # аналитика + декоратор для логирования
├─ go.sum
└─ go.mod
```
## Реализованные паттерны в проекте

Ниже — краткий список паттернов с указанием, где именно они применены и зачем.

---

### 1) **Factory** (Фабрика)
- **Где:** `domain/factory/*`  
- **Зачем:** Централизует создание доменных сущностей с валидацией (пустые ID/имя, тип категории ∈ {доход, расход}, `amount > 0`, валидная дата и т. п.).

---

### 2) **Facade** (Фасад)
- **Где:**  
  - `internal/importer/ImportFacade` — единая точка чтения из `DataImporter`, валидации через фабрики.  
  - `internal/export/ExportFacade` — единая точка экспорта всех сущностей (операции, категории, счета) в формате CSV/JSON.  
  - `internal/analitic/AnalyticsFacade` — обёртка над аналитическими операциями (расчёт баланса за период, группировки).
- **Зачем:** Прячет сложность связок (импорт/валидация/связывание/экспорт/аналитика) и предоставляет простой доступ.

---

### 3) **Strategy** (Стратегия)
- **Где:** `domain/ports.DataImporter` + реализации `internal/importer/{csv,json}`  
- **Зачем:** Внешний код (фасад импорта) не знает, из какого источника читать — подставляется нужная стратегия чтения (CSV/JSON).

---

### 4) **Visitor** (Посетитель)
- **Где:**  
  - Интерфейс: `internal/export/ExporterVisitor`  
  - Конкретные визиторы: `ExportCSV`, `ExportJSON`  
  - Элементы: адаптеры `CategoryAdapter`, `BankAccountAdapter`, `OperationAdapter` реализуют `Visitable.Accept(...)`
- **Зачем:** Легко добавлять новые форматы экспорта.

---

### 5) **Adapter** (Адаптер)
- **Где:** `internal/export/adaptersEntities.go`  
  - `CategoryAdapter`, `BankAccountAdapter`, `OperationAdapter`
- **Зачем:** у нас есть разные типы данных — счета, категории, операции. Экспортёры (CSV/JSON) ожидают единый интерфейс с методами `ExportBankAccount`, `ExportCategory`, `ExportOperation`.  
- **Проблема:** данные живут каждый в своём слайсе, а экспортёр хочет, чтобы их к нему подвели одинаковым способом.  
- **Решение (Adapter):** оборачиваем слайс в небольшой объект-адаптер с методом `Accept(exporter)`. Внутри он просто вызывает нужный метод экспортёра и передаёт ему свои данные.

---

### 6) **Decorator** (Декоратор)
- **Где:**  
  - Импорт: `LoggingImportDecorator` — логирует запуск/длительность/ошибки импорта.  
  - Экспорт: `LoginExportDecorator` — логирует запуск/длительность/ошибки экспорта по каждому адаптеру.  
  - Аналитика: `LoggingDecoratorForAnalitic` — логирует выполнение команд аналитики.
  - Аналитика: `InputPeriodDecorator` — добавляет даты для команды `CalculateBalanceCommand` так как она выводит данные за какой-то период
- **Зачем:** Добавляет к операциям сквозную функциональность (логирование) без изменения основной логики.

---

### 7) **Command** (Команда)
- **Где:** `internal/analitic` (интерфейс `Command` +  инициатор команд `Invoker`)
- **Зачем:** Операции аналитики оформлены как исполняемые команды; это упрощает логирование и позволдит в будущем добавлять еще команды.

---
### 8) **Шаблонный метод** 
- **Где:** `internal/importer/shablonFunc.go`
- Централизует **последовательность шагов** импорта.
- Позволяет **легко добавлять новые форматы** (XML, XLSX и т.п.): достаточно
  реализовать `ports.DataImporter`, а сценарий не меняется.



## Импорт:

**Поддерживаемые источники:** CSV и JSON через интерфейс `ports.DataImporter`  
**Фасад:** `ImportFacade`

- **Валидация через фабрики** (`domain/factory`):
  - `BankAccountFactory.CreateBankAccount` — проверяет `id`, `name`, `balance ≥ 0`.
  - `CategoryFactory.CreateCategory` — проверяет `id`, `name`, `type ∈ {доход, расход}`.
  - `OperationFactory.CreateOperation` — проверяет `id`, `type ∈ {доход, расход}`, `amount > 0`, `date`, ссылки.
- **Дубликаты ID**: сам импорт не сохраняет — дубликаты отсекаются на слое сохранения (`BankManager`).
- **Логи импорта**: фиксируются **время старта (RFC3339)**, **длительность (мс/с)**, источник (`CSV`/`JSON`).

**Формат CSV (ожидаемые колонки):**
- `bank_account`: `type=bank_account`, `id`→[1], `name`→[8], `balance`→[9]
- `category`: `type=category`, `id`→[1], `type(доход/расход)`→[2], `name`→[8]
- `operation`: `type=operation`, `id`→[1], `opType`→[2], `account_id`→[3], `amount`→[4], `date(YYYY-MM-DD)`→[5], `description`→[6], `category_id`→[7]

**Формат JSON:**
```json
{
  "bank_accounts": [{"id":"...","name":"...","balance":0}],
  "categories":    [{"id":"...","type":"доход|расход","name":"..."}],
  "operations":    [{"id":"...","type":"доход|расход","account_id":"...","amount":0,"date":"YYYY-MM-DD","description":"...","category_id":"..."}]
}
```
## Экспорт

Экспорт реализован через паттерн **Visitor**: конкретные экспортёры (`ExportCSV`, `ExportJSON`) посещают адаптеры данных (`BankAccountAdapter`, `CategoryAdapter`, `OperationAdapter`) и сохраняют их в файлы. Управляет всем **`ExportFacade`**, логирование выполняет `LoginExportDecorator`.

### Как работает
1. В TUI выбери: **Экспорт → CSV** или **Экспорт → JSON**.  
2. Введи базовое имя файла (например, `report`).  
3. Файлы будут сохранены в директории `./exportFiles` с суффиксами по типам данных.

### Интерфейсы и классы
- **ExporterVisitor** — интерфейс экспортёра:
  - `SetFilePath(path string) error`
  - `ExportBankAccount(accounts []*entities.BankAccount) error`
  - `ExportCategory(categories []*entities.Category) error`
  - `ExportOperation(operations []*entities.Operation) error`
- **Реализации**: `ExportCSV`, `ExportJSON`
- **Адаптеры (Visitable)**:
  - `BankAccountAdapter` — экспорт счетов
  - `CategoryAdapter` — экспорт категорий
  - `OperationAdapter` — экспорт операций
- **Фасад**: `ExportFacade` — конфигурирует адаптеры и запускает экспорт
- **Логи**: `LoginExportDecorator` — пишет в лог тип экспортёра, сущность, время старта и длительность

### Папка назначения
Все файлы экспортируются в директорию: exportFiles
Если её нет — будет создана автоматически. `SetFilePath("report")` нормализует имя и добавляет нужное расширение.

---

### Форматы файлов

#### CSV
Генерируются три файла:

- `report_BackAccounts.csv` — **счета**
  - Заголовки: `ID,Name,Balance`
- `report_Category.csv` — **категории**
  - Заголовки: `ID,Type,Name`
- `report_Operation.csv` — **операции**
  - Заголовки: `ID,Type,AccountID,Amount,Date,Description,CategoryID`

Пример `report_Operation.csv`:
```csv
ID,Type,AccountID,Amount,Date,Description,CategoryID
1,расход,1,1500.00,2024-01-15,Обед,2
2,доход,1,50000.00,2024-01-10,Зарплата,1
```
### JSON

Генерируются три файла:

- `report_BackAccounts.json`
- `report_Category.json`
- `report_Operation.json`

**Фрагмент `report_Operation.json`:**
```json
[
  {
    "id": "1",
    "type": "расход",
    "account": { "id": "1", "name": "Основной счёт", "balance": 0 },
    "amount": 1500,
    "date": "2024-01-15",
    "description": "Обед",
    "category": { "id": "2", "type": "расход", "name": "Еда" }
  }
]
```
## BankManager (кратко)

`BankManager` — простой слой над репозиториями.  
Через него приложение **добавляет / находит / удаляет / перечисляет** счета, категории и операции.


---

### Что делает

- Хранит ссылки на три репозитория:
  - `RepositoryBankAccounts`
  - `RepositoryCategory`
  - `RepositoryOperations`
- Массово добавляет данные.


---


// Баланс и ошибки
UpdateBalance(accountID string, newBalance float64) error
GetErrors() []errordata.ErrorRecord
