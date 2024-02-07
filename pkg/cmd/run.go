package cmd

import (
	"github.com/cidverse/go-vcsapp/pkg/task/taskcommon"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/cidverse/rewrite-app/pkg/openrewritelib"
	"github.com/cidverse/rewrite-app/pkg/tasks/openrewrite"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			dir, _ := cmd.Flags().GetString("dir")
			configFile, _ := cmd.Flags().GetString("config")
			if configFile == "" {
				configFile = "/app/config.yaml"
			}

			if dir == "" {
				runAsApp(configFile)
			} else {
				runLocal(dir, configFile)
			}
		},
	}
	cmd.Flags().Bool("dry-run", false, "Perform a dry run without making any changes")
	cmd.Flags().String("dir", "", "Directory of the project for local code generation")
	cmd.Flags().StringP("config", "c", "", "Configuration file")

	return cmd
}

func runAsApp(configFile string) {
	// tasks
	tasks := []taskcommon.Task{openrewrite.NewTask(configFile)}

	// platform
	platform, err := vcsapp.GetPlatformFromEnvironment()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to configure platform from environment")
	}

	// execute
	err = vcsapp.ExecuteTasks(platform, tasks)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to execute generate task")
	}
}

func runLocal(dir string, configFile string) {
	_, err := openrewritelib.RunOpenRewrite(dir, configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open-rewrite for local project")
	}
}
