# 🚀 Руководство по развертыванию Shuffle

## 📋 Предварительные требования

- Docker и Docker Compose
- Git
- Не менее 4GB RAM
- Не менее 10GB свободного места на диске

## 🔧 Быстрое развертывание

### 1. Клонирование репозитория
```bash
git clone https://github.com/Karamelishe/Shuffle.git
cd Shuffle
```

### 2. Проверка системы
```bash
# Запустите тест системы для проверки готовности
./test_system.sh
```

### 3. Запуск сервисов
```bash
# Запуск всех сервисов в фоне
docker compose up -d

# Проверка статуса сервисов
docker compose ps

# Просмотр логов
docker compose logs -f
```

### 4. Доступ к системе
- **Frontend**: http://localhost:3001
- **Backend API**: http://localhost:5001
- **OpenSearch**: http://localhost:9200

## 🔑 Система лицензирования

### Генерация лицензионного ключа
```bash
cd backend/license_generator

# Генерация базовой лицензии на 30 дней
go run standalone_license_generator.go generate basic 30

# Генерация профессиональной лицензии на 365 дней
go run standalone_license_generator.go generate professional 365

# Генерация корпоративной лицензии на 365 дней
go run standalone_license_generator.go generate enterprise 365
```

### Типы лицензий

| Тип | Пользователи | Workflow'ы | Выполнения | Особенности |
|-----|-------------|------------|------------|-------------|
| **Basic** | 3 | 5 | 1,000 | Базовый функционал |
| **Professional** | 10 | 50 | 10,000 | Расширенные интеграции |
| **Enterprise** | Без ограничений | Без ограничений | Без ограничений | Все функции + SSO |

### Активация лицензии
1. Сгенерируйте ключ с помощью генератора
2. Войдите в систему как администратор
3. Перейдите в настройки лицензии
4. Введите полученный ключ и активируйте

### Тестирование API лицензий
```bash
# Запуск тестов лицензионной системы
bash backend/test_license.sh
```

## 🌐 Интернационализация

Система поддерживает:
- 🇺🇸 **Английский** (по умолчанию)
- 🇷🇺 **Русский**

Переключение языка доступно в правом верхнем углу интерфейса.

## ⚙️ Конфигурация

### Основные переменные окружения (.env)
```bash
# Внешний hostname (измените на ваш домен или IP)
OUTER_HOSTNAME=localhost

# Порты для доступа
FRONTEND_PORT=3001
BACKEND_PORT=5001

# Пароли и безопасность
SHUFFLE_OPENSEARCH_PASSWORD=StrongShufflePassword321!
SHUFFLE_ENCRYPTION_MODIFIER=shuffle_default_encryption_key_change_me

# Учетная запись администратора по умолчанию
SHUFFLE_DEFAULT_USERNAME=admin
SHUFFLE_DEFAULT_PASSWORD=password
```

## 🛠️ Устранение неполадок

### Проблемы с запуском
```bash
# Проверка логов всех сервисов
docker compose logs

# Перезапуск сервисов
docker compose restart

# Полная пересборка
docker compose down
docker compose up -d --build
```

### Проблемы с лицензией
```bash
# Проверка статуса лицензии
curl -X GET http://localhost:5001/api/v1/license/info

# Валидация лицензионного ключа
cd backend/license_generator
go run standalone_license_generator.go validate YOUR_LICENSE_KEY
```

### Очистка данных
```bash
# Остановка и удаление всех данных
docker compose down -v
rm -rf /tmp/shuffle-*

# Пересоздание директорий
mkdir -p /tmp/shuffle-{database,apps,files}
```

## 📚 Дополнительные ресурсы

- **Документация по лицензированию**: `backend/LICENSE_SYSTEM.md`
- **README по лицензиям**: `README_LICENSE.md`
- **Настройка переводов**: `frontend/TRANSLATION_SETUP.md`

## 🆘 Поддержка

При возникновении проблем:
1. Проверьте логи: `docker compose logs`
2. Запустите тест системы: `./test_system.sh`
3. Проверьте состояние сервисов: `docker compose ps`
4. Убедитесь, что все порты свободны: `netstat -tlnp | grep -E '(3001|5001|9200)'`

## 🔄 Обновление системы

```bash
# Получение последних изменений
git pull origin main

# Пересборка и перезапуск
docker compose down
docker compose up -d --build
```