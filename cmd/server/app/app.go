package app

import (
	"net/url"
	"os"
	"path/filepath"
	"sort"

	"github.com/andrepinto/helmsman/api"
	"github.com/andrepinto/helmsman/pkg"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"k8s.io/helm/pkg/urlutil"
)

// NewCliApp ...
func NewCliApp() *cli.App {

	app := cli.NewApp()

	app.Name = "helmsman"
	app.Version = VERSION

	opts := NewHemlCmdOptionsCmdOptions()
	opts.AddFlags(app)

	app.Action = func(c *cli.Context) error {

		if opts.Debug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}

		var err error

		opts.Envs = c.StringSlice("env")

		opts.RepoUrl, err = urlutil.URLJoin(opts.RepoUrl, "/envs/%s/charts")
		if err != nil {
			log.Fatal(err)
		}

		opts.RepoUrl, _ = url.PathUnescape(opts.RepoUrl)

		err = Init(opts)
		if err != nil {
			log.Fatal(err)
		}

		proc := api.NewServer(&api.ServerOptions{
			Port:    opts.Port,
			RepoDir: opts.RepoDir,
			RepoUrl: opts.RepoUrl,
		})
		log.Debug(opts)
		return proc.Run()
	}

	// sort flags by name
	sort.Sort(cli.FlagsByName(app.Flags))

	return app
}

// Init ...
func Init(opts *HemlCmdOptions) error {

	log.Info(opts.Envs)

	for _, env := range opts.Envs {
		folder := filepath.Join(opts.RepoDir, env)
		err := os.MkdirAll(folder, 0777)
		if err != nil {
			return err
		}

		err = pkg.Index(opts.RepoDir, opts.RepoUrl, env, "")
		if err != nil {
			return err
		}
	}

	return nil

}
