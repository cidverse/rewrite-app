package openrewritelib

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/rs/zerolog/log"
)

func runGradleRewrite(projectDirectory string, config Config, rewriteYamlFile string) error {
	// generate gradle.init script
	gradleInitScript, err := RenderGradleInitTemplate(config, rewriteYamlFile)
	if err != nil {
		return fmt.Errorf("failed to render gradle.init script: %w", err)
	}

	gradleInitFile, err := os.CreateTemp("", "gradleinit")
	if err != nil {
		return fmt.Errorf("failed to create temp file for gradle.init script: %w", err)
	}
	defer os.Remove(gradleInitFile.Name())
	err = os.WriteFile(gradleInitFile.Name(), []byte(gradleInitScript), 0644)
	if err != nil {
		return fmt.Errorf("failed to write gradle.init script: %w", err)
	}

	// log the gradle.init script
	log.Trace().Str("template_id", "gradle.init").Msg(gradleInitScript)

	// ensure gradlew is executable
	err = os.Chmod(projectDirectory+"/gradlew", 0755)
	if err != nil {
		return fmt.Errorf("failed to add +x permission to gradlew: %w", err)
	}

	// run open-rewrite via gradle init script
	command := []string{
		projectDirectory + "/gradlew",
		"-Dorg.gradle.jvmargs=-Xmx4G",
		"rewriteRun",
		"--init-script=" + gradleInitFile.Name(),
		"-Drewrite.inPlace=true",
	}
	command = append(command, fmt.Sprintf("-Drewrite.configFile=%s", rewriteYamlFile))
	log.Debug().Strs("cmd", command).Msg("running command")

	cmd := exec.Command("bash", command...)
	cmd.Dir = projectDirectory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute openrewrite-gradle: %w", err)
	}

	return nil
}

// RenderGradleInitTemplate renders the gradle init template with the given config
func RenderGradleInitTemplate(config Config, rewriteYamlFile string) (string, error) {
	output, err := vcsapp.Render(getTemplate("gradle.init"), map[string]interface{}{
		"ConfigFile":   rewriteYamlFile,
		"Recipes":      config.Recipes,
		"Styles":       config.Styles,
		"ExcludeFiles": config.ExcludeFiles,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render description template: %w", err)
	}

	return string(output), nil
}
