# 🔐 Система лицензирования Shuffle

Полная система лицензирования для платформы автоматизации Shuffle с поддержкой русского языка.

## 📋 Что реализовано

### 🔧 Backend компоненты
- ✅ **license.go** - Основная логика лицензирования
- ✅ **standalone_license_generator.go** - CLI генератор лицензий
- ✅ **API endpoints** для управления лицензиями
- ✅ **Middleware** для защиты критических функций
- ✅ **Проверка лимитов** пользователей, workflow'ов и выполнений

### 🎨 Frontend компоненты
- ✅ **LicenseManager.jsx** - Интерфейс управления лицензиями
- ✅ **LicenseWarning.jsx** - Компонент предупреждений
- ✅ **Русская локализация** всех элементов интерфейса
- ✅ **Интеграция с существующей i18n системой**

### 📊 Типы лицензий

| Тип | Пользователи | Workflow'ы | Выполнения | Дополнительно |
|-----|-------------|------------|------------|---------------|
| **Basic** | 3 | 5 | 1,000/мес | Базовые приложения |
| **Professional** | 10 | 50 | 10,000/мес | Email поддержка, интеграции |
| **Enterprise** | ∞ | ∞ | ∞ | SSO, аудит, брендинг |

## 🚀 Быстрый старт

### 1. Генерация лицензии

```bash
cd backend/go-app
go run standalone_license_generator.go -generate -type professional -duration 365
```

**Пример вывода:**
```
License generated successfully!
================================
License Key:      79B93E6F-16D7B23D-1F465BFC-C931BD4A
Type:             professional
Expires:          2025-08-23 12:19:16
Max Users:        10
Max Workflows:    50
Max Executions:   10000
```

### 2. Валидация лицензии

```bash
go run standalone_license_generator.go -validate 79B93E6F-16D7B23D-1F465BFC-C931BD4A
```

### 3. Активация через API

```bash
curl -X POST http://localhost:5001/api/v1/license/activate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{"key": "79B93E6F-16D7B23D-1F465BFC-C931BD4A"}'
```

## 🔌 API Endpoints

### Активация лицензии
```http
POST /api/v1/license/activate
Content-Type: application/json
Authorization: Bearer YOUR_API_KEY

{
  "key": "XXXX-XXXX-XXXX-XXXX",
  "hardware_id": "optional"
}
```

### Информация о лицензии
```http
GET /api/v1/license/info
Authorization: Bearer YOUR_API_KEY
```

### Генерация лицензии (только администраторы)
```http
POST /api/v1/license/generate
Content-Type: application/json
Authorization: Bearer ADMIN_API_KEY

{
  "type": "professional",
  "organization_id": "org-123",
  "duration": 365
}
```

## 🛡️ Защищенные функции

Система автоматически защищает следующие операции:

- ✅ **Создание workflow'ов** - проверяет лимит workflow'ов
- ✅ **Выполнение workflow'ов** - проверяет лимит выполнений
- ✅ **Регистрация пользователей** - проверяет лимит пользователей
- ✅ **Доступ к расширенным функциям** - SSO, аудит и т.д.

## 🎨 Пользовательский интерфейс

### Компонент управления лицензиями
```jsx
import LicenseManager from './components/LicenseManager';

<LicenseManager
  globalUrl={globalUrl}
  userdata={userdata}
  theme={theme}
/>
```

### Предупреждения о лицензии
```jsx
import LicenseWarning from './components/LicenseWarning';

<LicenseWarning
  message="Требуется активная лицензия"
  onLicenseActivate={() => setShowLicenseManager(true)}
/>
```

## 🌍 Русская локализация

Все элементы интерфейса переведены на русский язык:

- 🔑 **Управление лицензиями** - полный интерфейс на русском
- ⚠️ **Предупреждения** - уведомления об истечении и лимитах
- 📊 **Статистика** - отображение лимитов и использования
- 🎯 **Типы лицензий** - описания планов на русском

## 🧪 Тестирование

### Автоматический тест API
```bash
chmod +x backend/test_license.sh
./backend/test_license.sh http://localhost:5001 YOUR_API_KEY
```

### Ручное тестирование
```bash
# 1. Генерируем лицензию
cd backend/go-app
go run standalone_license_generator.go -generate -type basic

# 2. Копируем ключ и активируем через UI
# 3. Проверяем ограничения при создании workflow'ов
```

## 📁 Структура файлов

```
backend/
├── go-app/
│   ├── license.go                      # Основная логика
│   ├── standalone_license_generator.go # CLI генератор
│   └── main.go                        # API endpoints + middleware
├── test_license.sh                    # Скрипт тестирования
└── LICENSE_SYSTEM.md                  # Подробная документация

frontend/
├── src/
│   ├── components/
│   │   ├── LicenseManager.jsx         # Управление лицензиями
│   │   └── LicenseWarning.jsx         # Предупреждения
│   └── i18n.js                        # Русские переводы
```

## ⚙️ Настройка и развертывание

### 1. Backend
```bash
# Файлы уже добавлены в main.go
# API endpoints автоматически доступны после запуска
```

### 2. Frontend
```bash
# Установить зависимость
cd frontend
npm install react-i18next@^12.0.0

# Компоненты готовы к использованию
```

### 3. База данных
- **Cloud**: Автоматически использует Google Datastore
- **On-premise**: Сохраняет в `/tmp/shuffle_license_*.json`

## 🔧 Интеграция с существующим кодом

### Проверка лицензии в коде
```go
// Проверка наличия лицензии
if !shuffle.IsLicensed(ctx, *org) {
    return errors.New("license required")
}

// Проверка конкретной функции
if !CheckLicenseFeature(ctx, orgID, "sso") {
    return errors.New("SSO not available")
}

// Проверка лимитов
limit := GetLicenseLimit(ctx, orgID, "workflows")
if limit > 0 && currentCount >= limit {
    return errors.New("workflow limit exceeded")
}
```

### Добавление middleware к новым endpoints
```go
r.HandleFunc("/api/v1/new-feature", 
    checkLicenseMiddleware(handleNewFeature, "advanced")).
    Methods("POST", "OPTIONS")
```

## 📈 Мониторинг и алерты

Система автоматически:
- 📊 Отслеживает использование лимитов
- ⏰ Предупреждает об истечении лицензии (за 30 дней)
- 🚫 Блокирует превышение лимитов
- 📝 Логирует все операции с лицензиями

## 🐛 Устранение неполадок

### Лицензия не активируется
```bash
# Проверьте ключ
go run standalone_license_generator.go -validate YOUR_KEY

# Проверьте API
curl -H "Authorization: Bearer API_KEY" http://localhost:5001/api/v1/license/info
```

### Middleware не работает
- Проверьте, что `license.go` подключен к main.go
- Убедитесь, что функция `shuffle.IsLicensed` доступна
- Проверьте логи на наличие ошибок аутентификации

### Frontend ошибки
- Убедитесь, что `react-i18next` установлен
- Проверьте, что компоненты правильно импортированы
- Проверьте консоль браузера на ошибки

## 🎯 Примеры использования

### Генерация корпоративной лицензии
```bash
go run standalone_license_generator.go \
  -generate \
  -type enterprise \
  -org "acme-corp" \
  -duration 730
```

### Массовая генерация лицензий
```bash
for i in {1..10}; do
  go run standalone_license_generator.go \
    -generate \
    -type basic \
    -org "client-$i" \
    -duration 365
done
```

### Проверка всех лицензий
```bash
go run standalone_license_generator.go -list
```

## 📚 Дополнительные ресурсы

- 📖 **Подробная документация**: `backend/LICENSE_SYSTEM.md`
- 🧪 **Скрипт тестирования**: `backend/test_license.sh`
- 🔧 **CLI помощь**: `go run standalone_license_generator.go -help`
- 🌐 **API документация**: Встроена в endpoints

## ✨ Заключение

Система лицензирования полностью готова к использованию и включает:

- 🔐 **Безопасную генерацию ключей** с криптографической защитой
- 🎨 **Красивый интерфейс** с полной русской локализацией
- 🛡️ **Автоматическую защиту** критических функций
- 📊 **Гибкие лимиты** для разных типов лицензий
- 🧪 **Полное тестирование** всех компонентов

Система готова для production использования! 🚀