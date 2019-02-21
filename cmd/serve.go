package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"net/http"

	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xescugc/notigator/assets"
	"github.com/xescugc/notigator/config"
	"github.com/xescugc/notigator/github"
	"github.com/xescugc/notigator/gitlab"
	"github.com/xescugc/notigator/service"
	"github.com/xescugc/notigator/trello"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the web server",
		Long:  "Starts the web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.New(viper.GetViper())
			if err != nil {
				return err
			}

			ctx := context.Background()

			gh := github.NewNotificationRepository(ctx, cfg.GitHubToken)
			gl := gitlab.NewNotificationRepository(ctx, cfg.GitLabToken)
			tr := trello.NewNotificationRepository(ctx, cfg.TrelloApiKey, cfg.TrelloToken, cfg.TrelloMember)

			s := service.New(gh, gl, tr)

			mux := http.NewServeMux()

			mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assets.AssetFS())))
			//mux.Handle("/assets/", http.FileServer(assets.AssetFS()))

			mux.Handle("/api/", service.MakeHandler(s))

			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				b, err := assets.Asset("templates/index.html")
				if err != nil {
					panic(err)
				}
				io.Copy(w, bytes.NewBuffer(b))
			})

			http.Handle("/", handlers.LoggingHandler(os.Stdout, mux))

			return http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
		},
	}
)

func init() {
	serveCmd.PersistentFlags().StringP("port", "p", "3000", "Destination port")
	viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))

	serveCmd.PersistentFlags().String("github-token", "", "Github Auth Token")
	viper.BindPFlag("github-token", serveCmd.PersistentFlags().Lookup("github-token"))

	serveCmd.PersistentFlags().String("gitlab-token", "", "Gitlab Auth Token")
	viper.BindPFlag("gitlab-token", serveCmd.PersistentFlags().Lookup("gitlab-token"))

	serveCmd.PersistentFlags().String("trello-token", "", "Trello Auth Token")
	viper.BindPFlag("trello-token", serveCmd.PersistentFlags().Lookup("trello-token"))

	serveCmd.PersistentFlags().String("trello-api-key", "", "Trello Api Key")
	viper.BindPFlag("trello-api-key", serveCmd.PersistentFlags().Lookup("trello-api-key"))

	serveCmd.PersistentFlags().String("trello-member", "", "Trello member")
	viper.BindPFlag("trello-member", serveCmd.PersistentFlags().Lookup("trello-member"))
}
