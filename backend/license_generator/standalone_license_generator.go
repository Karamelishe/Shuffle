package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// License структура для хранения информации о лицензии
type License struct {
	ID              string    `json:"id"`
	Key             string    `json:"key"`
	OrganizationID  string    `json:"organization_id"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	ActivatedAt     *time.Time `json:"activated_at,omitempty"`
	MaxUsers        int       `json:"max_users"`
	MaxWorkflows    int       `json:"max_workflows"`
	MaxExecutions   int       `json:"max_executions"`
	Features        []string  `json:"features"`
	HardwareID      string    `json:"hardware_id"`
	LastValidated   time.Time `json:"last_validated"`
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

// generateID генерирует уникальный ID
func generateID() string {
	return uuid.NewV4().String()
}

// generateLicenseKey генерирует уникальный лицензионный ключ
func generateLicenseKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	hash := sha256.Sum256(bytes)
	key := hex.EncodeToString(hash[:])[:32]
	
	formatted := fmt.Sprintf("%s-%s-%s-%s", 
		strings.ToUpper(key[0:8]), 
		strings.ToUpper(key[8:16]), 
		strings.ToUpper(key[16:24]), 
		strings.ToUpper(key[24:32]))
	
	return formatted, nil
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
		license.MaxUsers = -1
		license.MaxWorkflows = -1
		license.MaxExecutions = -1
	}

	// Сохраняем лицензию в файл
	data, err := json.MarshalIndent(license, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal license: %v", err)
	}
	
	filename := fmt.Sprintf("/tmp/shuffle_license_%s.json", license.ID)
	err = ioutil.WriteFile(filename, data, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to save license: %v", err)
	}

	return license, nil
}

// GetLicenseByKey получает лицензию по ключу (простая реализация для standalone)
func GetLicenseByKey(ctx context.Context, key string) (*License, error) {
	// В standalone режиме ищем во всех файлах лицензий
	files, err := ioutil.ReadDir("/tmp")
	if err != nil {
		return nil, fmt.Errorf("failed to read license directory: %v", err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "shuffle_license_") && strings.HasSuffix(file.Name(), ".json") {
			data, err := ioutil.ReadFile("/tmp/" + file.Name())
			if err != nil {
				continue
			}

			var license License
			if err := json.Unmarshal(data, &license); err != nil {
				continue
			}

			if license.Key == key {
				return &license, nil
			}
		}
	}

	return nil, fmt.Errorf("license not found")
}

// ValidateLicense проверяет валидность лицензии
func ValidateLicense(ctx context.Context, license *License) error {
	if license == nil {
		return fmt.Errorf("license is nil")
	}

	if license.Status != "active" {
		return fmt.Errorf("license is not active")
	}

	if time.Now().After(license.ExpiresAt) {
		return fmt.Errorf("license has expired")
	}

	return nil
}

func main() {
	var (
		licenseType    = flag.String("type", "basic", "License type: basic, professional, enterprise")
		organizationID = flag.String("org", "", "Organization ID (optional)")
		duration       = flag.Int("duration", 365, "License duration in days")
		showHelp       = flag.Bool("help", false, "Show help")
		generate       = flag.Bool("generate", false, "Generate a new license")
		validate       = flag.String("validate", "", "Validate a license key")
		list           = flag.Bool("list", false, "List all licenses")
	)

	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	ctx := context.Background()

	if *generate {
		generateLicense(ctx, *licenseType, *organizationID, *duration)
	} else if *validate != "" {
		validateLicenseKey(ctx, *validate)
	} else if *list {
		listLicenses(ctx)
	} else {
		fmt.Println("Please specify an action. Use -help for more information.")
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Shuffle License Generator (Standalone)")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run standalone_license_generator.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -generate            Generate a new license")
	fmt.Println("  -type string         License type: basic, professional, enterprise (default: basic)")
	fmt.Println("  -org string          Organization ID (optional)")
	fmt.Println("  -duration int        License duration in days (default: 365)")
	fmt.Println("  -validate string     Validate a license key")
	fmt.Println("  -list               List all licenses")
	fmt.Println("  -help               Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Generate a basic license for 1 year")
	fmt.Println("  go run standalone_license_generator.go -generate -type basic -duration 365")
	fmt.Println()
	fmt.Println("  # Generate a professional license for specific organization")
	fmt.Println("  go run standalone_license_generator.go -generate -type professional -org org-123 -duration 730")
	fmt.Println()
	fmt.Println("  # Validate a license key")
	fmt.Println("  go run standalone_license_generator.go -validate ABCD1234-EFGH5678-IJKL9012-MNOP3456")
	fmt.Println()
	fmt.Println("  # List all licenses")
	fmt.Println("  go run standalone_license_generator.go -list")
	fmt.Println()
	fmt.Println("License Types:")
	fmt.Println("  basic         - Up to 3 users, 5 workflows, 1000 executions/month")
	fmt.Println("  professional  - Up to 10 users, 50 workflows, 10000 executions/month")
	fmt.Println("  enterprise    - Unlimited users, workflows, and executions")
	fmt.Println()
	fmt.Println("Note: This is a standalone version. Licenses are stored in /tmp/ directory.")
}

func generateLicense(ctx context.Context, licenseType, organizationID string, duration int) {
	fmt.Printf("Generating %s license for %d days...\n", licenseType, duration)

	if _, exists := LicenseFeatures[licenseType]; !exists {
		fmt.Printf("Error: Invalid license type '%s'. Valid types: basic, professional, enterprise\n", licenseType)
		os.Exit(1)
	}

	durationTime := time.Duration(duration) * 24 * time.Hour
	license, err := CreateLicense(ctx, licenseType, organizationID, durationTime)
	if err != nil {
		fmt.Printf("Error generating license: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("License generated successfully!")
	fmt.Println("================================")
	fmt.Printf("License Key:      %s\n", license.Key)
	fmt.Printf("License ID:       %s\n", license.ID)
	fmt.Printf("Type:             %s\n", license.Type)
	fmt.Printf("Status:           %s\n", license.Status)
	fmt.Printf("Created:          %s\n", license.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Expires:          %s\n", license.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Max Users:        %s\n", formatLimit(license.MaxUsers))
	fmt.Printf("Max Workflows:    %s\n", formatLimit(license.MaxWorkflows))
	fmt.Printf("Max Executions:   %s\n", formatLimit(license.MaxExecutions))
	if organizationID != "" {
		fmt.Printf("Organization ID:  %s\n", organizationID)
	}
	fmt.Println()
	fmt.Println("Features:")
	for _, feature := range license.Features {
		fmt.Printf("  - %s\n", feature)
	}
	fmt.Println()
	fmt.Printf("License saved to: /tmp/shuffle_license_%s.json\n", license.ID)
	fmt.Println("IMPORTANT: Save this license key in a secure location!")
}

func validateLicenseKey(ctx context.Context, key string) {
	fmt.Printf("Validating license key: %s\n", key)

	license, err := GetLicenseByKey(ctx, key)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = ValidateLicense(ctx, license)
	if err != nil {
		fmt.Printf("License validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("License is valid!")
	fmt.Println("=================")
	fmt.Printf("License ID:       %s\n", license.ID)
	fmt.Printf("Type:             %s\n", license.Type)
	fmt.Printf("Status:           %s\n", license.Status)
	fmt.Printf("Created:          %s\n", license.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Expires:          %s\n", license.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Max Users:        %s\n", formatLimit(license.MaxUsers))
	fmt.Printf("Max Workflows:    %s\n", formatLimit(license.MaxWorkflows))
	fmt.Printf("Max Executions:   %s\n", formatLimit(license.MaxExecutions))

	if license.OrganizationID != "" {
		fmt.Printf("Organization ID:  %s\n", license.OrganizationID)
	}

	if license.HardwareID != "" {
		fmt.Printf("Hardware ID:      %s\n", license.HardwareID)
	}

	if license.ActivatedAt != nil {
		fmt.Printf("Activated:        %s\n", license.ActivatedAt.Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("Last Validated:   %s\n", license.LastValidated.Format("2006-01-02 15:04:05"))

	timeLeft := time.Until(license.ExpiresAt)
	if timeLeft > 0 {
		days := int(timeLeft.Hours() / 24)
		fmt.Printf("Time remaining:   %d days\n", days)
	} else {
		fmt.Println("Status:           EXPIRED")
	}

	fmt.Println()
	fmt.Println("Features:")
	for _, feature := range license.Features {
		fmt.Printf("  - %s\n", feature)
	}
}

func listLicenses(ctx context.Context) {
	fmt.Println("Listing all licenses...")
	
	files, err := ioutil.ReadDir("/tmp")
	if err != nil {
		fmt.Printf("Error reading license directory: %v\n", err)
		return
	}

	licenseCount := 0
	fmt.Println("Found licenses:")
	fmt.Println("===============")

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "shuffle_license_") && strings.HasSuffix(file.Name(), ".json") {
			data, err := ioutil.ReadFile("/tmp/" + file.Name())
			if err != nil {
				continue
			}

			var license License
			if err := json.Unmarshal(data, &license); err != nil {
				continue
			}

			licenseCount++
			fmt.Printf("%d. %s (%s)\n", licenseCount, license.Key, license.Type)
			fmt.Printf("   Status: %s, Expires: %s\n", license.Status, license.ExpiresAt.Format("2006-01-02"))
			if license.OrganizationID != "" {
				fmt.Printf("   Organization: %s\n", license.OrganizationID)
			}
			fmt.Println()
		}
	}

	if licenseCount == 0 {
		fmt.Println("No licenses found in /tmp/ directory.")
	} else {
		fmt.Printf("Total: %d licenses\n", licenseCount)
	}
}

func formatLimit(limit int) string {
	if limit == -1 {
		return "unlimited"
	}
	return strconv.Itoa(limit)
}