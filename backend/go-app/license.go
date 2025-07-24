package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	uuid "github.com/satori/go.uuid"
	"github.com/shuffle/shuffle-shared"
)

// License структура для хранения информации о лицензии
type License struct {
	ID              string    `json:"id" datastore:"id"`
	Key             string    `json:"key" datastore:"key"`
	OrganizationID  string    `json:"organization_id" datastore:"organization_id"`
	Type            string    `json:"type" datastore:"type"` // "basic", "professional", "enterprise"
	Status          string    `json:"status" datastore:"status"` // "active", "expired", "revoked"
	CreatedAt       time.Time `json:"created_at" datastore:"created_at"`
	ExpiresAt       time.Time `json:"expires_at" datastore:"expires_at"`
	ActivatedAt     *time.Time `json:"activated_at,omitempty" datastore:"activated_at"`
	MaxUsers        int       `json:"max_users" datastore:"max_users"`
	MaxWorkflows    int       `json:"max_workflows" datastore:"max_workflows"`
	MaxExecutions   int       `json:"max_executions" datastore:"max_executions"`
	Features        []string  `json:"features" datastore:"features"`
	HardwareID      string    `json:"hardware_id" datastore:"hardware_id"`
	LastValidated   time.Time `json:"last_validated" datastore:"last_validated"`
}

// LicenseRequest структура для запроса активации лицензии
type LicenseRequest struct {
	Key        string `json:"key"`
	HardwareID string `json:"hardware_id"`
}

// LicenseResponse структура для ответа на запрос лицензии
type LicenseResponse struct {
	Success   bool    `json:"success"`
	Message   string  `json:"message"`
	License   *License `json:"license,omitempty"`
}

// LicenseFeatures определяет доступные функции для разных типов лицензий
var LicenseFeatures = map[string][]string{
	"basic": {
		"workflows:5",
		"users:3",
		"executions:1000",
		"apps:basic",
		"support:community",
	},
	"professional": {
		"workflows:50",
		"users:10",
		"executions:10000",
		"apps:all",
		"support:email",
		"integrations:advanced",
		"reporting:basic",
	},
	"enterprise": {
		"workflows:unlimited",
		"users:unlimited",
		"executions:unlimited",
		"apps:all",
		"support:priority",
		"integrations:all",
		"reporting:advanced",
		"sso:enabled",
		"audit:enabled",
		"custom_branding:enabled",
	},
}

// generateLicenseKey генерирует уникальный лицензионный ключ
func generateLicenseKey() (string, error) {
	// Генерируем 32 байта случайных данных
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Создаем хэш SHA256
	hash := sha256.Sum256(bytes)
	
	// Конвертируем в hex и берем первые 32 символа
	key := hex.EncodeToString(hash[:])[:32]
	
	// Форматируем ключ в группы по 8 символов, разделенные дефисами
	formatted := fmt.Sprintf("%s-%s-%s-%s", 
		strings.ToUpper(key[0:8]), 
		strings.ToUpper(key[8:16]), 
		strings.ToUpper(key[16:24]), 
		strings.ToUpper(key[24:32]))
	
	return formatted, nil
}

// generateID генерирует уникальный ID
func generateID() string {
	return uuid.NewV4().String()
}

// generateHardwareID генерирует уникальный идентификатор оборудования
func generateHardwareID() string {
	// В реальной реализации здесь должна быть логика получения
	// уникальных характеристик оборудования (MAC адрес, серийный номер процессора и т.д.)
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// writeFile записывает данные в файл
func writeFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0600)
}

// readFile читает данные из файла
func readFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// CreateLicense создает новую лицензию
func CreateLicense(ctx context.Context, licenseType string, organizationID string, duration time.Duration) (*License, error) {
	key, err := generateLicenseKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate license key: %v", err)
	}

	features, exists := LicenseFeatures[licenseType]
	if !exists {
		return nil, fmt.Errorf("unknown license type: %s", licenseType)
	}

	license := &License{
		ID:             generateID(),
		Key:            key,
		OrganizationID: organizationID,
		Type:           licenseType,
		Status:         "active",
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(duration),
		Features:       features,
		LastValidated:  time.Now(),
	}

	// Устанавливаем лимиты на основе типа лицензии
	switch licenseType {
	case "basic":
		license.MaxUsers = 3
		license.MaxWorkflows = 5
		license.MaxExecutions = 1000
	case "professional":
		license.MaxUsers = 10
		license.MaxWorkflows = 50
		license.MaxExecutions = 10000
	case "enterprise":
		license.MaxUsers = -1 // unlimited
		license.MaxWorkflows = -1 // unlimited
		license.MaxExecutions = -1 // unlimited
	}

	// Сохраняем лицензию в базе данных
	err = SetLicense(ctx, *license)
	if err != nil {
		return nil, fmt.Errorf("failed to save license: %v", err)
	}

	return license, nil
}

// SetLicense сохраняет лицензию в базе данных
func SetLicense(ctx context.Context, license License) error {
	if runningEnvironment == "cloud" {
		dbclient := shuffle.GetDatastore()
		key := datastore.NameKey("License", license.ID, nil)
		if _, err := dbclient.Put(ctx, key, &license); err != nil {
			return err
		}
	} else {
		// Для локальной установки сохраняем в файл
		data, err := json.Marshal(license)
		if err != nil {
			return err
		}
		
		filename := fmt.Sprintf("/tmp/shuffle_license_%s.json", license.ID)
		return writeFile(filename, data)
	}
	
	return nil
}

// GetLicense получает лицензию по ID
func GetLicense(ctx context.Context, licenseID string) (*License, error) {
	if runningEnvironment == "cloud" {
		dbclient := shuffle.GetDatastore()
		key := datastore.NameKey("License", licenseID, nil)
		license := &License{}
		if err := dbclient.Get(ctx, key, license); err != nil {
			return nil, err
		}
		return license, nil
	} else {
		// Для локальной установки читаем из файла
		filename := fmt.Sprintf("/tmp/shuffle_license_%s.json", licenseID)
		data, err := readFile(filename)
		if err != nil {
			return nil, err
		}
		
		license := &License{}
		err = json.Unmarshal(data, license)
		if err != nil {
			return nil, err
		}
		return license, nil
	}
}

// GetLicenseByKey получает лицензию по ключу
func GetLicenseByKey(ctx context.Context, key string) (*License, error) {
	if runningEnvironment == "cloud" {
		dbclient := shuffle.GetDatastore()
		query := datastore.NewQuery("License").Filter("key =", key).Limit(1)
		licenses := []License{}
		_, err := dbclient.GetAll(ctx, query, &licenses)
		if err != nil {
			return nil, err
		}
		
		if len(licenses) == 0 {
			return nil, fmt.Errorf("license not found")
		}
		
		return &licenses[0], nil
	} else {
		// Для локальной установки проверяем все файлы лицензий
		// В реальной реализации лучше использовать индекс или базу данных
		return nil, fmt.Errorf("local license lookup by key not implemented")
	}
}

// GetOrganizationLicense получает активную лицензию организации
func GetOrganizationLicense(ctx context.Context, organizationID string) (*License, error) {
	if runningEnvironment == "cloud" {
		dbclient := shuffle.GetDatastore()
		query := datastore.NewQuery("License").
			Filter("organization_id =", organizationID).
			Filter("status =", "active").
			Order("-created_at").
			Limit(1)
		
		licenses := []License{}
		_, err := dbclient.GetAll(ctx, query, &licenses)
		if err != nil {
			return nil, err
		}
		
		if len(licenses) == 0 {
			return nil, fmt.Errorf("no active license found for organization")
		}
		
		return &licenses[0], nil
	} else {
		// Для локальной установки ищем в файлах
		return nil, fmt.Errorf("local organization license lookup not implemented")
	}
}

// ActivateLicense активирует лицензию
func ActivateLicense(ctx context.Context, key string, hardwareID string, organizationID string) (*License, error) {
	license, err := GetLicenseByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("license not found: %v", err)
	}

	// Проверяем статус лицензии
	if license.Status != "active" {
		return nil, fmt.Errorf("license is not active")
	}

	// Проверяем срок действия
	if time.Now().After(license.ExpiresAt) {
		license.Status = "expired"
		SetLicense(ctx, *license)
		return nil, fmt.Errorf("license has expired")
	}

	// Проверяем, не активирована ли уже лицензия на другом оборудовании
	if license.HardwareID != "" && license.HardwareID != hardwareID {
		return nil, fmt.Errorf("license is already activated on different hardware")
	}

	// Активируем лицензию
	now := time.Now()
	license.ActivatedAt = &now
	license.HardwareID = hardwareID
	license.OrganizationID = organizationID
	license.LastValidated = now

	err = SetLicense(ctx, *license)
	if err != nil {
		return nil, fmt.Errorf("failed to update license: %v", err)
	}

	return license, nil
}

// ValidateLicense проверяет валидность лицензии
func ValidateLicense(ctx context.Context, license *License) error {
	if license == nil {
		return fmt.Errorf("license is nil")
	}

	// Проверяем статус
	if license.Status != "active" {
		return fmt.Errorf("license is not active")
	}

	// Проверяем срок действия
	if time.Now().After(license.ExpiresAt) {
		license.Status = "expired"
		SetLicense(ctx, *license)
		return fmt.Errorf("license has expired")
	}

	// Обновляем время последней проверки
	license.LastValidated = time.Now()
	SetLicense(ctx, *license)

	return nil
}

// checkLicenseMiddleware проверяет лицензию перед выполнением функции
func checkLicenseMiddleware(next http.HandlerFunc, requiredFeature string) http.HandlerFunc {
	return func(resp http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		user, err := handleApiAuthentication(resp, request)
		if err != nil {
			log.Printf("[WARNING] Api authentication failed in license check: %s", err)
			resp.WriteHeader(401)
			resp.Write([]byte(`{"success": false, "reason": "Failed authentication"}`))
			return
		}

		// Получаем организацию
		org, err := shuffle.GetOrg(ctx, user.ActiveOrg.Id)
		if err != nil {
			log.Printf("[WARNING] Failed to get organization: %v", err)
			resp.WriteHeader(500)
			resp.Write([]byte(`{"success": false, "reason": "Failed to get organization"}`))
			return
		}

		// Проверяем лицензию
		if !shuffle.IsLicensed(ctx, *org) {
			resp.WriteHeader(403)
			resp.Write([]byte(`{"success": false, "reason": "License required. Please activate a valid license to access this feature.", "license_required": true}`))
			return
		}

		// Если указана конкретная функция, проверяем её доступность
		if requiredFeature != "" {
			if !CheckLicenseFeature(ctx, user.ActiveOrg.Id, requiredFeature) {
				resp.WriteHeader(403)
				resp.Write([]byte(fmt.Sprintf(`{"success": false, "reason": "Feature '%s' not available in your license plan", "license_required": true}`, requiredFeature)))
				return
			}
		}

		// Проверяем лимиты
		if err := checkLicenseLimits(ctx, user.ActiveOrg.Id); err != nil {
			resp.WriteHeader(403)
			resp.Write([]byte(fmt.Sprintf(`{"success": false, "reason": "%s", "license_required": true}`, err.Error())))
			return
		}

		// Продолжаем выполнение
		next(resp, request)
	}
}

// checkLicenseLimits проверяет лимиты лицензии
func checkLicenseLimits(ctx context.Context, organizationID string) error {
	license, err := GetOrganizationLicense(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("failed to get license: %v", err)
	}

	// Проверяем лимит пользователей
	if license.MaxUsers > 0 {
		users, err := shuffle.GetUsersByOrg(ctx, organizationID)
		if err == nil && len(users) >= license.MaxUsers {
			return fmt.Errorf("user limit reached (%d/%d)", len(users), license.MaxUsers)
		}
	}

	// Проверяем лимит workflow'ов
	if license.MaxWorkflows > 0 {
		workflows, err := shuffle.GetWorkflows(ctx, 1000, organizationID)
		if err == nil && len(workflows) >= license.MaxWorkflows {
			return fmt.Errorf("workflow limit reached (%d/%d)", len(workflows), license.MaxWorkflows)
		}
	}

	return nil
}

// CheckLicenseFeature проверяет, доступна ли определенная функция в лицензии
func CheckLicenseFeature(ctx context.Context, organizationID string, feature string) bool {
	license, err := GetOrganizationLicense(ctx, organizationID)
	if err != nil {
		return false
	}

	err = ValidateLicense(ctx, license)
	if err != nil {
		return false
	}

	for _, f := range license.Features {
		if strings.HasPrefix(f, feature) {
			return true
		}
	}

	return false
}

// GetLicenseLimit получает лимит для определенного ресурса
func GetLicenseLimit(ctx context.Context, organizationID string, resource string) int {
	license, err := GetOrganizationLicense(ctx, organizationID)
	if err != nil {
		return 0
	}

	err = ValidateLicense(ctx, license)
	if err != nil {
		return 0
	}

	switch resource {
	case "users":
		return license.MaxUsers
	case "workflows":
		return license.MaxWorkflows
	case "executions":
		return license.MaxExecutions
	default:
		return 0
	}
}

// handleGenerateLicense обрабатывает запрос на генерацию новой лицензии (только для администраторов)
func handleGenerateLicense(resp http.ResponseWriter, request *http.Request) {
	cors := handleCors(resp, request)
	if cors {
		return
	}

	ctx := context.Background()
	user, err := handleApiAuthentication(resp, request)
	if err != nil {
		log.Printf("[WARNING] Api authentication failed in generate license: %s", err)
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false, "reason": "Failed authentication"}`))
		return
	}

	// Проверяем права администратора
	if user.Role != "admin" {
		resp.WriteHeader(403)
		resp.Write([]byte(`{"success": false, "reason": "Admin access required"}`))
		return
	}

	body, err := readBody(request)
	if err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(`{"success": false, "reason": "Failed to read body"}`))
		return
	}

	type GenerateLicenseRequest struct {
		Type           string `json:"type"`
		OrganizationID string `json:"organization_id"`
		Duration       int    `json:"duration"` // в днях
	}

	var req GenerateLicenseRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(`{"success": false, "reason": "Failed to parse JSON"}`))
		return
	}

	// Валидация
	if req.Type == "" {
		req.Type = "basic"
	}
	if req.Duration <= 0 {
		req.Duration = 365 // по умолчанию 1 год
	}

	duration := time.Duration(req.Duration) * 24 * time.Hour
	license, err := CreateLicense(ctx, req.Type, req.OrganizationID, duration)
	if err != nil {
		log.Printf("[ERROR] Failed to create license: %v", err)
		resp.WriteHeader(500)
		resp.Write([]byte(`{"success": false, "reason": "Failed to create license"}`))
		return
	}

	response := LicenseResponse{
		Success: true,
		Message: "License generated successfully",
		License: license,
	}

	respData, err := json.Marshal(response)
	if err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte(`{"success": false, "reason": "Failed to marshal response"}`))
		return
	}

	resp.WriteHeader(200)
	resp.Write(respData)
}

// handleActivateLicense обрабатывает запрос на активацию лицензии
func handleActivateLicense(resp http.ResponseWriter, request *http.Request) {
	cors := handleCors(resp, request)
	if cors {
		return
	}

	ctx := context.Background()
	user, err := handleApiAuthentication(resp, request)
	if err != nil {
		log.Printf("[WARNING] Api authentication failed in activate license: %s", err)
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false, "reason": "Failed authentication"}`))
		return
	}

	body, err := readBody(request)
	if err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(`{"success": false, "reason": "Failed to read body"}`))
		return
	}

	var req LicenseRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(`{"success": false, "reason": "Failed to parse JSON"}`))
		return
	}

	if req.Key == "" {
		resp.WriteHeader(400)
		resp.Write([]byte(`{"success": false, "reason": "License key is required"}`))
		return
	}

	// Генерируем hardware ID если не предоставлен
	if req.HardwareID == "" {
		req.HardwareID = generateHardwareID()
	}

	license, err := ActivateLicense(ctx, req.Key, req.HardwareID, user.ActiveOrg.Id)
	if err != nil {
		log.Printf("[ERROR] Failed to activate license: %v", err)
		resp.WriteHeader(400)
		resp.Write([]byte(fmt.Sprintf(`{"success": false, "reason": "%s"}`, err.Error())))
		return
	}

	response := LicenseResponse{
		Success: true,
		Message: "License activated successfully",
		License: license,
	}

	respData, err := json.Marshal(response)
	if err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte(`{"success": false, "reason": "Failed to marshal response"}`))
		return
	}

	resp.WriteHeader(200)
	resp.Write(respData)
}

// handleGetLicenseInfo обрабатывает запрос на получение информации о лицензии
func handleGetLicenseInfo(resp http.ResponseWriter, request *http.Request) {
	cors := handleCors(resp, request)
	if cors {
		return
	}

	ctx := context.Background()
	user, err := handleApiAuthentication(resp, request)
	if err != nil {
		log.Printf("[WARNING] Api authentication failed in get license info: %s", err)
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false, "reason": "Failed authentication"}`))
		return
	}

	license, err := GetOrganizationLicense(ctx, user.ActiveOrg.Id)
	if err != nil {
		resp.WriteHeader(404)
		resp.Write([]byte(`{"success": false, "reason": "No license found"}`))
		return
	}

	err = ValidateLicense(ctx, license)
	if err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(fmt.Sprintf(`{"success": false, "reason": "%s"}`, err.Error())))
		return
	}

	response := LicenseResponse{
		Success: true,
		Message: "License information retrieved successfully",
		License: license,
	}

	respData, err := json.Marshal(response)
	if err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte(`{"success": false, "reason": "Failed to marshal response"}`))
		return
	}

	resp.WriteHeader(200)
	resp.Write(respData)
}