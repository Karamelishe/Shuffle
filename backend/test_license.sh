#!/bin/bash

# Скрипт для тестирования API лицензирования Shuffle
# Использование: ./test_license.sh [base_url] [api_key]

BASE_URL=${1:-"http://localhost:5001"}
API_KEY=${2:-"your-api-key-here"}

echo "🔑 Тестирование системы лицензирования Shuffle"
echo "================================================"
echo "Base URL: $BASE_URL"
echo "API Key: $API_KEY"
echo ""

# Функция для выполнения HTTP запросов
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo "📡 $description"
    echo "   $method $endpoint"
    
    if [ -n "$data" ]; then
        echo "   Data: $data"
        response=$(curl -s -X $method \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $API_KEY" \
            -d "$data" \
            "$BASE_URL$endpoint")
    else
        response=$(curl -s -X $method \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $API_KEY" \
            "$BASE_URL$endpoint")
    fi
    
    echo "   Response: $response"
    echo ""
    
    # Возвращаем response для дальнейшего использования
    echo "$response"
}

# 1. Проверяем текущую лицензию
echo "1️⃣ Проверка текущей лицензии"
license_info=$(make_request "GET" "/api/v1/license/info" "" "Получение информации о лицензии")

# Проверяем, есть ли лицензия
if echo "$license_info" | grep -q '"success":true'; then
    echo "✅ Лицензия найдена"
    # Извлекаем ключ лицензии для дальнейших тестов
    license_key=$(echo "$license_info" | grep -o '"key":"[^"]*"' | cut -d'"' -f4)
    echo "   Ключ лицензии: $license_key"
else
    echo "⚠️  Активная лицензия не найдена"
    license_key=""
fi

echo ""

# 2. Генерируем новую лицензию (требует права администратора)
echo "2️⃣ Генерация новой лицензии"
generate_data='{
    "type": "professional",
    "organization_id": "test-org-123",
    "duration": 365
}'

new_license=$(make_request "POST" "/api/v1/license/generate" "$generate_data" "Генерация новой лицензии")

if echo "$new_license" | grep -q '"success":true'; then
    echo "✅ Лицензия успешно сгенерирована"
    # Извлекаем новый ключ
    new_license_key=$(echo "$new_license" | grep -o '"key":"[^"]*"' | cut -d'"' -f4)
    echo "   Новый ключ: $new_license_key"
else
    echo "❌ Не удалось сгенерировать лицензию (возможно, нужны права администратора)"
    new_license_key=""
fi

echo ""

# 3. Активируем лицензию (если сгенерировали новую)
if [ -n "$new_license_key" ]; then
    echo "3️⃣ Активация новой лицензии"
    activate_data="{
        \"key\": \"$new_license_key\",
        \"hardware_id\": \"test-hardware-$(date +%s)\"
    }"
    
    activation_result=$(make_request "POST" "/api/v1/license/activate" "$activate_data" "Активация лицензии")
    
    if echo "$activation_result" | grep -q '"success":true'; then
        echo "✅ Лицензия успешно активирована"
    else
        echo "❌ Не удалось активировать лицензию"
    fi
else
    echo "3️⃣ Пропускаем активацию (нет нового ключа)"
fi

echo ""

# 4. Тестируем защищенные endpoints
echo "4️⃣ Тестирование защищенных endpoints"

# Тест создания workflow (требует лицензию)
echo "   Тестируем создание workflow..."
workflow_data='{
    "name": "Test License Workflow",
    "description": "Workflow для тестирования лицензии",
    "tags": ["test", "license"]
}'

workflow_result=$(make_request "POST" "/api/v1/workflows" "$workflow_data" "Создание workflow")

if echo "$workflow_result" | grep -q '"success":true'; then
    echo "✅ Workflow создан успешно"
elif echo "$workflow_result" | grep -q 'license_required'; then
    echo "⚠️  Создание workflow заблокировано - требуется лицензия"
else
    echo "❌ Ошибка при создании workflow"
fi

echo ""

# 5. Проверяем лимиты
echo "5️⃣ Проверка лимитов лицензии"

# Получаем информацию о пользователях
users_result=$(make_request "GET" "/api/v1/users" "" "Получение списка пользователей")

if echo "$users_result" | grep -q '"success":true'; then
    echo "✅ Список пользователей получен"
    # Можно добавить подсчет пользователей
elif echo "$users_result" | grep -q 'license_required'; then
    echo "⚠️  Доступ к пользователям заблокирован - требуется лицензия"
else
    echo "❌ Ошибка при получении пользователей"
fi

echo ""

# 6. Финальная проверка лицензии
echo "6️⃣ Финальная проверка лицензии"
final_license_info=$(make_request "GET" "/api/v1/license/info" "" "Финальная проверка лицензии")

if echo "$final_license_info" | grep -q '"success":true'; then
    echo "✅ Лицензия активна"
    
    # Извлекаем и отображаем основную информацию
    echo ""
    echo "📊 Информация о лицензии:"
    echo "$final_license_info" | python3 -m json.tool 2>/dev/null || echo "$final_license_info"
else
    echo "❌ Лицензия не активна или недоступна"
fi

echo ""
echo "🏁 Тестирование завершено"
echo "========================"

# Проверяем, установлен ли jq для красивого вывода JSON
if command -v jq &> /dev/null; then
    echo ""
    echo "💡 Для лучшего отображения JSON используйте jq:"
    echo "   curl -H \"Authorization: Bearer $API_KEY\" $BASE_URL/api/v1/license/info | jq"
fi

echo ""
echo "📚 Документация: backend/LICENSE_SYSTEM.md"
echo "🔧 CLI инструмент: cd backend/go-app && go run license_generator.go -help"