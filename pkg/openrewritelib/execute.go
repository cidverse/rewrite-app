package openrewritelib

import (
	_ "embed"
	"fmt"
	"os"
	"path"

	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/rs/zerolog/log"
)

func RunOpenRewrite(projectDirectory string, configFile string) (Config, error) {
	cfg, configErr := LoadProjectConfig(configFile)
	if configErr != nil {
		return Config{}, configErr
	}

	// rewrite.yml config
	rewriteYamlFile := path.Join(projectDirectory, "rewrite.yml")

	// if file is missing
	if _, err := os.Stat(rewriteYamlFile); os.IsNotExist(err) {
		tempFile, err := os.CreateTemp("", "rewrite")
		if err != nil {
			return Config{}, fmt.Errorf("failed to create temp file for rewrite.yml: %w", err)
		}
		defer os.Remove(tempFile.Name())

		// render
		output, err := RenderRewriteConfigTemplate(cfg)
		if err != nil {
			return Config{}, fmt.Errorf("failed to render rewrite.yml: %w", err)
		}

		// write template to temp file
		err = os.WriteFile(tempFile.Name(), output, 0644)
		if err != nil {
			return Config{}, fmt.Errorf("failed to write rewrite.yml: %w", err)
		}

		rewriteYamlFile = tempFile.Name()

		// log the rewrite.yml
		log.Trace().Str("template_id", "rewrite.yml").Msg(string(output))
	}

	// run
	if _, err := os.Stat(path.Join(projectDirectory, "build.gradle.kts")); err == nil {
		return cfg, runGradleRewrite(projectDirectory, cfg, rewriteYamlFile)
	} else if _, err := os.Stat(path.Join(projectDirectory, "build.gradle")); err == nil {
		return cfg, runGradleRewrite(projectDirectory, cfg, rewriteYamlFile)
	} else if _, err := os.Stat(path.Join(projectDirectory, "pom.xml")); err == nil {
		return cfg, fmt.Errorf("maven is not supported yet")
	} else {
		return cfg, UnsupportedProjectError
	}
}

// RenderRewriteConfigTemplate renders the rewrite config template
func RenderRewriteConfigTemplate(config Config) ([]byte, error) {
	output, err := vcsapp.Render(getTemplate("rewrite.yml"), map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to render description template: %w", err)
	}

	return output, nil
}

// getEnvOrDefault retrieves the value of the environment variable with the given name.
func getEnvOrDefault(name, defaultValue string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return defaultValue
}
