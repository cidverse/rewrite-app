package openrewritelib

import (
	_ "embed"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

//go:embed templates/gradleinit.gohtml
var gradleInitTemplate []byte

//go:embed templates/rewrite.yml
var rewriteYamlTemplate []byte

func getTemplate(id string) string {
	// user-provided templates
	if len(os.Getenv("REWRITEAPP_TEMPLATE_DIR")) > 0 {
		if _, err := os.Stat(path.Join(os.Getenv("REWRITEAPP_TEMPLATE_DIR"), id)); err == nil {
			bytes, readErr := os.ReadFile(path.Join(os.Getenv("REWRITEAPP_TEMPLATES"), id))
			if readErr != nil {
				log.Fatal().Err(readErr).Str("id", id).Msg("failed to read template")
			}
			return string(bytes)
		}
	}

	// embedded templates
	if id == "gradle.init" {
		return string(gradleInitTemplate)
	} else if id == "rewrite.yml" {
		return string(rewriteYamlTemplate)
	}

	log.Fatal().Str("id", id).Msg("unknown template id")
	return ""
}
