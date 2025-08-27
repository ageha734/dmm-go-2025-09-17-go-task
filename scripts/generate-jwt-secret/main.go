package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	defaultLength        = 64
	minLength            = 16
	recommendedMinLength = 32
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

type Config struct {
	Length     int
	UpdateEnv  bool
	Force      bool
	OutputOnly bool
	Help       bool
}

type JWTSecretGenerator struct {
	config Config
}

func NewJWTSecretGenerator(config Config) *JWTSecretGenerator {
	return &JWTSecretGenerator{config: config}
}

func printColored(color, prefix, message string) {
	if !isOutputOnly() {
		fmt.Printf("%s%s %s%s\n", color, prefix, message, colorReset)
	}
}

func printInfo(message string) {
	printColored(colorBlue, "ℹ️ ", message)
}

func printSuccess(message string) {
	printColored(colorGreen, "✅", message)
}

func printWarning(message string) {
	printColored(colorYellow, "⚠️ ", message)
}

func printError(message string) {
	printColored(colorRed, "❌", message)
}

var outputOnly bool

func isOutputOnly() bool {
	return outputOnly
}

func showUsage() {
	fmt.Println("JWT Secret Generator")
	fmt.Println("")
	fmt.Println("Usage: go run cmd/generate-jwt-secret/main.go [options]")
	fmt.Println("   or: ./generate-jwt-secret [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -length <number>      Set the length of the secret (default: 64)")
	fmt.Println("  -update-env          Update .env file with the generated secret")
	fmt.Println("  -force               Force update even if JWT_SECRET already exists")
	fmt.Println("  -output-only         Only output the secret without any messages")
	fmt.Println("  -help                Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run scripts/generate-jwt-secret/main.go                    # Generate and display a 64-character secret")
	fmt.Println("  go run scripts/generate-jwt-secret/main.go -length 128        # Generate a 128-character secret")
	fmt.Println("  go run scripts/generate-jwt-secret/main.go -update-env        # Generate and update .env file")
	fmt.Println("  go run scripts/generate-jwt-secret/main.go -output-only       # Only output the secret (useful for scripts)")
	fmt.Println("")
}

func (g *JWTSecretGenerator) generateSecret() (string, error) {
	byteLength := (g.config.Length * 3) / 4
	if (g.config.Length*3)%4 != 0 {
		byteLength++
	}

	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	encoded := base64.URLEncoding.EncodeToString(randomBytes)

	encoded = strings.TrimRight(encoded, "=")
	if len(encoded) > g.config.Length {
		encoded = encoded[:g.config.Length]
	}

	for len(encoded) < g.config.Length {
		additionalBytes := make([]byte, 1)
		if _, err := rand.Read(additionalBytes); err != nil {
			return "", fmt.Errorf("failed to generate additional random bytes: %w", err)
		}
		additionalChar := base64.URLEncoding.EncodeToString(additionalBytes)
		additionalChar = strings.TrimRight(additionalChar, "=")
		if len(additionalChar) > 0 {
			encoded += string(additionalChar[0])
		}
	}

	return encoded[:g.config.Length], nil
}

func (g *JWTSecretGenerator) validateSecret(secret string) error {
	length := len(secret)

	if length < minLength {
		return fmt.Errorf("secret is too short (%d characters). Minimum length is %d", length, minLength)
	}

	if length < recommendedMinLength {
		printWarning(fmt.Sprintf("Secret is shorter than recommended (%d characters). Recommended minimum is %d.", length, recommendedMinLength))
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range secret {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	diversityScore := 0
	if hasUpper {
		diversityScore++
	}
	if hasLower {
		diversityScore++
	}
	if hasDigit {
		diversityScore++
	}
	if hasSpecial {
		diversityScore++
	}

	if diversityScore < 2 {
		printWarning("Secret has low character diversity. Consider regenerating.")
	}

	return nil
}

func findProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}

func validateFilePath(path string) error {
	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected: %s", path)
	}

	if filepath.IsAbs(cleanPath) {
		return nil
	}

	return nil
}

func (g *JWTSecretGenerator) updateEnvFile(secret string) error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	envFile := filepath.Join(projectRoot, ".env")
	envTemplate := filepath.Join(projectRoot, ".env.template")

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		printInfo(".env file not found. Creating from template...")
		if _, err := os.Stat(envTemplate); err == nil {
			if err := copyFile(envTemplate, envFile); err != nil {
				return fmt.Errorf("failed to copy template: %w", err)
			}
			printSuccess(".env file created from template")
		} else {
			return fmt.Errorf("neither .env nor .env.template found")
		}
	}

	if !g.config.Force {
		if currentSecret, exists, err := getCurrentJWTSecret(envFile); err != nil {
			return fmt.Errorf("failed to check current JWT_SECRET: %w", err)
		} else if exists && currentSecret != "your-secret-key-here" && currentSecret != "" {
			printWarning(fmt.Sprintf("JWT_SECRET already exists in .env file: %s", currentSecret))
			printWarning("Use -force to overwrite it")
			return fmt.Errorf("JWT_SECRET already exists")
		}
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFile := fmt.Sprintf("%s.backup.%s", envFile, timestamp)
	if err := copyFile(envFile, backupFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	printInfo(fmt.Sprintf("Created backup: %s", backupFile))

	if err := updateJWTSecretInFile(envFile, secret); err != nil {
		return fmt.Errorf("failed to update .env file: %w", err)
	}

	printSuccess("JWT_SECRET updated in .env file")
	return nil
}

func copyFile(src, dst string) error {
	if err := validateFilePath(src); err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}
	if err := validateFilePath(dst); err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}

	sourceFile, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := sourceFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close source file: %v\n", closeErr)
		}
	}()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close destination file: %v\n", closeErr)
		}
	}()

	scanner := bufio.NewScanner(sourceFile)
	writer := bufio.NewWriter(destFile)

	for scanner.Scan() {
		if _, err := writer.WriteString(scanner.Text() + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func getCurrentJWTSecret(envFile string) (string, bool, error) {
	if err := validateFilePath(envFile); err != nil {
		return "", false, fmt.Errorf("invalid env file path: %w", err)
	}

	file, err := os.Open(filepath.Clean(envFile))
	if err != nil {
		return "", false, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
		}
	}()

	scanner := bufio.NewScanner(file)
	jwtSecretRegex := regexp.MustCompile(`^JWT_SECRET=(.*)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := jwtSecretRegex.FindStringSubmatch(line); matches != nil {
			return matches[1], true, nil
		}
	}

	return "", false, scanner.Err()
}

func updateJWTSecretInFile(envFile, secret string) error {
	if err := validateFilePath(envFile); err != nil {
		return fmt.Errorf("invalid env file path: %w", err)
	}

	file, err := os.Open(filepath.Clean(envFile))
	if err != nil {
		return err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	jwtSecretRegex := regexp.MustCompile(`^JWT_SECRET=.*$`)
	jwtSecretUpdated := false

	for scanner.Scan() {
		line := scanner.Text()
		if jwtSecretRegex.MatchString(line) {
			lines = append(lines, fmt.Sprintf("JWT_SECRET=%s", secret))
			jwtSecretUpdated = true
		} else {
			lines = append(lines, line)
		}
	}
	if closeErr := file.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !jwtSecretUpdated {
		lines = append(lines, fmt.Sprintf("JWT_SECRET=%s", secret))
	}

	outputFile, err := os.Create(envFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close output file: %v\n", closeErr)
		}
	}()

	writer := bufio.NewWriter(outputFile)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func parseFlags() Config {
	var config Config

	flag.IntVar(&config.Length, "length", defaultLength, "Set the length of the secret")
	flag.BoolVar(&config.UpdateEnv, "update-env", false, "Update .env file with the generated secret")
	flag.BoolVar(&config.Force, "force", false, "Force update even if JWT_SECRET already exists")
	flag.BoolVar(&config.OutputOnly, "output-only", false, "Only output the secret without any messages")
	flag.BoolVar(&config.Help, "help", false, "Show help message")

	flag.Parse()

	outputOnly = config.OutputOnly

	return config
}

func (g *JWTSecretGenerator) Run() error {
	if g.config.Help {
		showUsage()
		return nil
	}

	if g.config.Length < minLength {
		return fmt.Errorf("length must be >= %d", minLength)
	}

	secret, err := g.generateSecret()
	if err != nil {
		return fmt.Errorf("failed to generate secret: %w", err)
	}

	if err := g.validateSecret(secret); err != nil {
		return fmt.Errorf("secret validation failed: %w", err)
	}

	if g.config.OutputOnly {
		fmt.Println(secret)
	} else {
		printInfo(fmt.Sprintf("Generated JWT Secret (%d characters):", len(secret)))
		fmt.Println(secret)
		fmt.Println()
	}

	if g.config.UpdateEnv {
		if err := g.updateEnvFile(secret); err != nil {
			return err
		}

		if !g.config.OutputOnly {
			fmt.Println()
			printSuccess("JWT Secret generation completed!")
			printInfo("Your .env file has been updated with the new JWT_SECRET")
			printWarning("Make sure to restart your application to use the new secret")
		}
	} else if !g.config.OutputOnly {
		fmt.Println()
		printInfo("To update your .env file, run:")
		fmt.Printf("  go run cmd/generate-jwt-secret/main.go -update-env\n")
		fmt.Println()
		printInfo("Or copy the secret above and manually update your .env file:")
		fmt.Printf("  JWT_SECRET=%s\n", secret)
	}

	return nil
}

func main() {
	config := parseFlags()
	generator := NewJWTSecretGenerator(config)

	if err := generator.Run(); err != nil {
		printError(err.Error())
		os.Exit(1)
	}
}
