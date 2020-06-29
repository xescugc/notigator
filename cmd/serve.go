package cmd

import (
	"bytes"
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
	"github.com/xescugc/notigator/immem"
	"github.com/xescugc/notigator/notification"
	"github.com/xescugc/notigator/service"
	"github.com/xescugc/notigator/source"
	"github.com/xescugc/notigator/trello"
	"github.com/xescugc/notigator/zeplin"
)

func init() {
	serveCmd.PersistentFlags().StringP("port", "p", "3000", "Destination port")
	viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))
}

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

			notifications := make(map[string]notification.Repository)
			sources := make([]source.Source, 0, len(cfg.Sources))

			// Maps all the sources from the configuration to source.Source
			// and initializes the notification repositories needed
			for _, s := range cfg.Sources {
				srcCan, err := source.CanonicalString(s.Canonical)
				if err != nil {
					return fmt.Errorf("invalid canonical on config %q: %s", s.Canonical, err)
				}

				src := source.Source{
					Name:      s.Name,
					Canonical: srcCan,
				}
				src.BuildID()
				sources = append(sources, src)

				nr, err := initializeRepository(srcCan, s)
				if err != nil {
					return fmt.Errorf("could not initializa repository: %s", err)
				}

				if _, ok := notifications[src.ID]; ok {
					return fmt.Errorf("there is already a repeated Name+Canonical combination on the confi %q", src.ID)
				}

				notifications[src.ID] = nr
			}

			srcr := immem.NewSourceRepository(sources)

			s := service.New(srcr, notifications)

			mux := http.NewServeMux()

			mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assets.AssetFS())))

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

func initializeRepository(sc source.Canonical, cs config.Source) (notification.Repository, error) {
	switch sc {
	case source.Github:
		return github.NewNotificationRepository(cs.Token), nil
	case source.Gitlab:
		r, err := gitlab.NewNotificationRepository(cs.Token)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize GitLab source: %w", err)
		}
		return r, nil
	case source.Trello:
		return trello.NewNotificationRepository(cs.APIKey, cs.Token), nil
	case source.Zeplin:
		return zeplin.NewNotificationRepository(cs.Token), nil
	default:
		return nil, fmt.Errorf("not implemented source %q", sc)
	}
}
