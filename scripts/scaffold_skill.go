package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./scripts/scaffold_skill.go <skill_name>")
		os.Exit(1)
	}

	name := normalizeName(os.Args[1])
	if name == "" {
		fmt.Println("Error: invalid skill name")
		os.Exit(1)
	}
	title := titleCase(name)

	interfacePath := filepath.Join("core", "skills", name+".go")
	pluginDir := filepath.Join("core", "plugins", name)
	pluginPath := filepath.Join(pluginDir, "plugin.go")

	if exists(interfacePath) || exists(pluginPath) {
		fmt.Printf("Error: skill '%s' already exists\n", name)
		os.Exit(1)
	}

	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		panic(err)
	}

	interfaceCode := fmt.Sprintf(`package skills

// %s defines the capability contract for the %s skill.
type %s interface {
	// Execute performs the primary skill action.
	Execute(param string) error
}
`, title, name, title)

	pluginCode := fmt.Sprintf(`package %s

import "fmt"

// Plugin is the default implementation for %s.
type Plugin struct{}

func New() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Execute(param string) error {
	fmt.Printf("[%s Plugin] Executing: %%s\n", param)
	return nil
}
`, name, name, title)

	writeFile(interfacePath, interfaceCode)
	writeFile(pluginPath, pluginCode)

	fmt.Printf("Created skill: %s\n", name)
	fmt.Printf("1. Interface: %s\n", interfacePath)
	fmt.Printf("2. Plugin:    %s\n", pluginPath)
	fmt.Println("Next: register it in core/engine/registry.go")
}

func normalizeName(in string) string {
	out := strings.ToLower(strings.TrimSpace(in))
	out = strings.ReplaceAll(out, "-", "_")
	return out
}

func titleCase(in string) string {
	parts := strings.Split(in, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		panic(err)
	}
}
