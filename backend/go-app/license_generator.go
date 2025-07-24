package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// LicenseGenerator CLI инструмент для генерации лицензий
func main() {
	// Определяем флаги командной строки
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
	fmt.Println("Shuffle License Generator")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run license_generator.go [options]")
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
	fmt.Println("  go run license_generator.go -generate -type basic -duration 365")
	fmt.Println()
	fmt.Println("  # Generate a professional license for specific organization")
	fmt.Println("  go run license_generator.go -generate -type professional -org org-123 -duration 730")
	fmt.Println()
	fmt.Println("  # Validate a license key")
	fmt.Println("  go run license_generator.go -validate ABCD1234-EFGH5678-IJKL9012-MNOP3456")
	fmt.Println()
	fmt.Println("  # List all licenses")
	fmt.Println("  go run license_generator.go -list")
	fmt.Println()
	fmt.Println("License Types:")
	fmt.Println("  basic         - Up to 3 users, 5 workflows, 1000 executions/month")
	fmt.Println("  professional  - Up to 10 users, 50 workflows, 10000 executions/month")
	fmt.Println("  enterprise    - Unlimited users, workflows, and executions")
}

func generateLicense(ctx context.Context, licenseType, organizationID string, duration int) {
	fmt.Printf("Generating %s license for %d days...\n", licenseType, duration)

	// Проверяем валидность типа лицензии
	if _, exists := LicenseFeatures[licenseType]; !exists {
		fmt.Printf("Error: Invalid license type '%s'. Valid types: basic, professional, enterprise\n", licenseType)
		os.Exit(1)
	}

	// Генерируем лицензию
	durationTime := time.Duration(duration) * 24 * time.Hour
	license, err := CreateLicense(ctx, licenseType, organizationID, durationTime)
	if err != nil {
		fmt.Printf("Error generating license: %v\n", err)
		os.Exit(1)
	}

	// Выводим информацию о созданной лицензии
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

	// Проверяем сколько времени осталось
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
	fmt.Println("Note: This feature requires database access and may not work in standalone mode.")
	
	// В реальной реализации здесь была бы логика получения всех лицензий
	// Для простоты выводим заглушку
	fmt.Println("License listing not implemented for standalone mode.")
	fmt.Println("Use the web interface or API to list licenses.")
}

func formatLimit(limit int) string {
	if limit == -1 {
		return "unlimited"
	}
	return strconv.Itoa(limit)
}

// Простая версия для standalone запуска
func init() {
	// Инициализируем переменные окружения для standalone режима
	if os.Getenv("SHUFFLE_LICENSE_STANDALONE") == "true" {
		runningEnvironment = "onprem"
	}
}