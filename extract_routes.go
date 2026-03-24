package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	baseDir := `.\modules`

	re := regexp.MustCompile(`\.(GET|POST|PUT|PATCH|DELETE)\("([^"]+)"`)

	pathsMap := make(map[string]map[string]string)

	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".go" && strings.HasSuffix(path, "routes.go") {
			content, _ := os.ReadFile(path)
			lines := strings.Split(string(content), "\n")
			module := filepath.Base(filepath.Dir(path))

			for _, line := range lines {
				matches := re.FindAllStringSubmatch(line, -1)
				for _, match := range matches {
					if len(match) > 2 {
						method := strings.ToLower(match[1])
						route := match[2]

						// Replace Gin /:id with OpenAPI /{id}
						route = regexp.MustCompile(`:([a-zA-Z0-9_]+)`).ReplaceAllString(route, "{$1}")

						if pathsMap[route] == nil {
							pathsMap[route] = make(map[string]string)
						}
						pathsMap[route][method] = module
					}
				}
			}
		}
		return nil
	})

	existingContent, err := os.ReadFile(`c:\Users\Drako\Desktop\cloud-tecno\aura-back\infrastructure\docs\api\openapi.yaml`)
	if err != nil {
		fmt.Printf("Error reading existing yaml: %v\n", err)
		return
	}
	existingStr := string(existingContent)

	var sb strings.Builder

	// Create map of existing routes
	existingRoutesMap := make(map[string]bool)
	existingLines := strings.Split(existingStr, "\n")
	for _, l := range existingLines {
		if strings.HasPrefix(l, "  /") {
			routeText := strings.TrimRight(strings.TrimSpace(l), ":")
			existingRoutesMap[routeText] = true
		}
	}

	for path, methods := range pathsMap {
		if !existingRoutesMap[path] {
			sb.WriteString(fmt.Sprintf("  %s:\n", path))
			for method, module := range methods {
				sb.WriteString(fmt.Sprintf("    %s:\n", method))

				moduleTitle := strings.Title(strings.ReplaceAll(module, "-", " "))

				sb.WriteString(fmt.Sprintf("      summary: %s operation for %s\n", strings.ToUpper(method), moduleTitle))
				sb.WriteString(fmt.Sprintf("      tags:\n        - %s\n", moduleTitle))
				sb.WriteString("      security:\n        - bearerAuth: []\n")

				params := regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`).FindAllStringSubmatch(path, -1)
				if len(params) > 0 {
					sb.WriteString("      parameters:\n")
					for _, param := range params {
						sb.WriteString(fmt.Sprintf("        - in: path\n          name: %s\n          required: true\n          schema:\n            type: string\n", param[1]))
					}
				}

				sb.WriteString("      responses:\n        \"200\":\n          description: Successful response\n        \"401\":\n          description: Unauthorized\n")
			}
		}
	}

	filePath := `C:\Users\Drako\Desktop\cloud-tecno\aura-back\new_routes.yaml`
	err = os.WriteFile(filePath, []byte(sb.String()), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	} else {
		fmt.Println("Generated", filePath)
	}
}
