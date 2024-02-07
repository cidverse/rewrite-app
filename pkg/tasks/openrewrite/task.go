package openrewrite

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/cidverse/go-vcsapp/pkg/task/simpletask"
	"github.com/cidverse/go-vcsapp/pkg/task/taskcommon"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/cidverse/rewrite-app/pkg/openrewritelib"
	"github.com/rs/zerolog/log"
)

const branchName = "refactor/open-rewrite"

//go:embed templates/description.gohtml
var descriptionTemplate []byte

type OpenRewriteTask struct {
	configFile string
}

// Name returns the name of the task
func (n OpenRewriteTask) Name() string {
	return "generate"
}

// Execute runs the task
func (n OpenRewriteTask) Execute(ctx taskcommon.TaskContext) error {
	log.Info().Str("dir", ctx.Directory).Str("project", ctx.Repository.Namespace+"/"+ctx.Repository.Name).Msg("running open-rewrite task")

	// create helper
	helper := simpletask.New(ctx)

	// clone repository
	err := helper.Clone()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// create and checkout new branch
	err = helper.CreateBranch(branchName)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	// run open rewrite
	cfg, err := openrewritelib.RunOpenRewrite(ctx.Directory, n.configFile)
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	// commit message and description
	commitMessage := fmt.Sprintf("refactor: run open-rewrite recipes")
	description, err := vcsapp.Render(string(descriptionTemplate), map[string]interface{}{
		"PlatformName": ctx.Platform.Name(),
		"PlatformSlug": ctx.Platform.Slug(),
		"Recipes":      cfg.Recipes,
		"Styles":       cfg.Styles,
		"Footer":       os.Getenv("REWRITEAPP_FOOTER_HIDE") != "true",
		"FooterCustom": os.Getenv("REWRITEAPP_FOOTER_CUSTOM"),
	})
	if err != nil {
		return fmt.Errorf("failed to render description template: %w", err)
	}

	// commit push and create or update merge request
	err = helper.CommitPushAndMergeRequest(commitMessage, commitMessage, string(description), "")
	if err != nil {
		return fmt.Errorf("failed to commit push and create or update merge request: %w", err)
	}

	return nil
}

func NewTask(configFile string) OpenRewriteTask {
	return OpenRewriteTask{
		configFile: configFile,
	}
}
