package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

type OutputFile struct {
	TemplatePath string
	OutputPath   string
}

func Generate(cfg *Config, projectRoot string) error {
	files := buildOutputFiles(cfg, projectRoot)

	tmpl := template.New("gen").Funcs(template.FuncMap{
		"title":        title,
		"lower":        strings.ToLower,
		"upper":        strings.ToUpper,
		"snake":        toSnakeCase,
		"hasPrefix":    strings.HasPrefix,
		"trimPrefix":   strings.TrimPrefix,
		"fieldGoType":  fieldGoType,
		"statusTypeName": statusTypeName,
	})

	entries, err := templateFS.ReadDir("templates")
	if err != nil {
		return fmt.Errorf("read templates dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := templateFS.ReadFile("templates/" + entry.Name())
		if err != nil {
			return fmt.Errorf("read template %s: %w", entry.Name(), err)
		}
		if _, err := tmpl.New(entry.Name()).Parse(string(data)); err != nil {
			return fmt.Errorf("parse template %s: %w", entry.Name(), err)
		}
	}

	var generated, skipped []string
	for _, f := range files {
		tmplName := f.TemplatePath
		t := tmpl.Lookup(tmplName)
		if t == nil {
			return fmt.Errorf("template %s not found", tmplName)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, cfg); err != nil {
			return fmt.Errorf("execute template %s: %w", tmplName, err)
		}

		if _, err := os.Stat(f.OutputPath); err == nil {
			skipped = append(skipped, f.OutputPath)
			continue
		}

		dir := filepath.Dir(f.OutputPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create dir %s: %w", dir, err)
		}

		if err := os.WriteFile(f.OutputPath, buf.Bytes(), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", f.OutputPath, err)
		}
		generated = append(generated, f.OutputPath)
	}

	fmt.Println("\n=== Generated Files ===")
	for _, f := range generated {
		fmt.Printf("  ✓ %s\n", f)
	}
	if len(skipped) > 0 {
		fmt.Println("\n=== Skipped (already exist) ===")
		for _, f := range skipped {
			fmt.Printf("  ✗ %s\n", f)
		}
	}

	printInstructions(cfg, projectRoot)
	return nil
}

func buildOutputFiles(cfg *Config, root string) []OutputFile {
	m := cfg.Module
	g := cfg.Group
	l := cfg.LabelEn

	serverRoot := filepath.Join(root, "server")
	adminRoot := filepath.Join(root, "admin")
	moduleDir := filepath.Join(serverRoot, "internal", "modules", g, m)
	modelDir := filepath.Join(serverRoot, "internal", "platform", "model")
	frontModuleDir := filepath.Join(adminRoot, "src", "modules", g)

	files := []OutputFile{
		// Backend
		{TemplatePath: "model.go.tpl", OutputPath: filepath.Join(modelDir, m+".go")},
		{TemplatePath: "domain_types.go.tpl", OutputPath: filepath.Join(moduleDir, "domain", "types.go")},
		{TemplatePath: "api_dto.go.tpl", OutputPath: filepath.Join(moduleDir, "api", "dto.go")},
		{TemplatePath: "api_routes.go.tpl", OutputPath: filepath.Join(moduleDir, "api", "routes.go")},
		{TemplatePath: "api_handler.go.tpl", OutputPath: filepath.Join(moduleDir, "api", "handler.go")},
		{TemplatePath: "app_ports.go.tpl", OutputPath: filepath.Join(moduleDir, "application", "ports.go")},
		{TemplatePath: "app_service.go.tpl", OutputPath: filepath.Join(moduleDir, "application", "service.go")},
		{TemplatePath: "infra_repo.go.tpl", OutputPath: filepath.Join(moduleDir, "infra", "repository.go")},
		{TemplatePath: "routes.go.tpl", OutputPath: filepath.Join(moduleDir, "routes.go")},
		{TemplatePath: "services.go.tpl", OutputPath: filepath.Join(moduleDir, "services.go")},
		// Frontend
		{TemplatePath: "types_ts.tpl", OutputPath: filepath.Join(frontModuleDir, "types", m+".ts")},
		{TemplatePath: "api_ts.tpl", OutputPath: filepath.Join(frontModuleDir, "api", m+".ts")},
		{TemplatePath: "composable_ts.tpl", OutputPath: filepath.Join(frontModuleDir, "composables", "use"+l+"Page.ts")},
		{TemplatePath: "utils_ts.tpl", OutputPath: filepath.Join(frontModuleDir, "composables", m+"-page.utils.ts")},
		{TemplatePath: "form_modal_vue.tpl", OutputPath: filepath.Join(frontModuleDir, "components", l+"FormModal.vue")},
		{TemplatePath: "filter_bar_vue.tpl", OutputPath: filepath.Join(frontModuleDir, "components", l+"FilterBar.vue")},
		{TemplatePath: "table_vue.tpl", OutputPath: filepath.Join(frontModuleDir, "components", l+"Table.vue")},
		{TemplatePath: "view_vue.tpl", OutputPath: filepath.Join(frontModuleDir, "pages", l+"View.vue")},
	}

	return files
}

func printInstructions(cfg *Config, root string) {
	m := cfg.Module
	g := cfg.Group
	l := cfg.LabelEn

	fmt.Printf("\n=== Manual Steps Required ===\n")
	fmt.Printf("\n1. Register route in server/internal/modules/%s/routes.go:\n", g)
	fmt.Printf("   %smodule.RegisterRoutes(system, %smodule.RouteOptions{\n", m, m)
	fmt.Printf("       DB:  opts.DB,\n")
	fmt.Printf("       Log: opts.Log,\n")
	fmt.Printf("   })\n")
	fmt.Printf("   Import: %smodule \"ez-admin-gin/server/internal/modules/%s/%s\"\n", m, g, m)

	fmt.Printf("\n2. Create database table (migration):\n")
	fmt.Printf("   Table: %s\n", cfg.Table)

	fmt.Printf("\n3. Add menu icon in admin/src/router/dynamic-menu.ts:\n")
	fmt.Printf("   '%s:%s': <icon>,\n", g, m)

	fmt.Printf("\n4. Add menu/permission seed data:\n")
	fmt.Printf("   Menu code: %s:%s\n", g, m)
	fmt.Printf("   Permissions: %s:%s:list, %s:%s:create, %s:%s:update, %s:%s:delete\n", g, m, g, m, g, m, g, m)
	if cfg.HasStatus {
		fmt.Printf("   Status: %s:%s:update_status\n", g, m)
	}

	fmt.Printf("\n5. Frontend route is auto-discovered via:\n")
	fmt.Printf("   admin/src/modules/%s/pages/%sView.vue\n", g, l)
}

func toSnakeCase(s string) string {
	var result []byte
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, byte(c+32))
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}

func fieldGoType(f Field) string {
	if strings.HasPrefix(f.Type, "model.") {
		return f.Type
	}
	return f.Type
}

func statusTypeName(f Field) string {
	if strings.HasPrefix(f.Type, "model.") {
		return strings.TrimPrefix(f.Type, "model.")
	}
	return "int"
}
