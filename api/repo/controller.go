package repo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrepinto/helmsman/pkg"
	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

const MimeGzip = "application/gzip"

func (pr *RepoResource) chartCtrl(request *restful.Request, response *restful.Response) {

	log.Debug("get chart")

	id := request.PathParameter("chart")
	env := request.PathParameter("env")

	response.Header().Set("Content-Type", "text/plain")

	file, err := os.Open(filepath.Join(pr.RepoDir, env, id))
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, "500: Charts error.")
		return
	}

	b, err := io.ReadAll(file)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, "500: Charts error.")
		return
	}

	accepts := strings.Split(request.Request.Header.Get("Accept"), ",")
	for _, accept := range accepts {
		if accept == MimeGzip {
			response.Header().Set("Content-Type", MimeGzip)
			break
		}
	}

	response.Write(b)
}

func (pr *RepoResource) uploadChartCtrl(request *restful.Request, response *restful.Response) {

	log.Debug("upload a chart")

	id := request.PathParameter("chart")
	env := request.PathParameter("env")

	response.Header().Set("Content-Type", "text/plain")

	f, err := os.Create(filepath.Join(pr.RepoDir, env, id))
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, "500: Charts error.")
		return
	}

	defer f.Close()

	_, err = io.Copy(f, request.Request.Body)
	defer request.Request.Body.Close()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, "500: Charts error.")
		return
	}

	err = pkg.Index(pr.RepoDir, pr.RepoUrl, env, "")
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("500: Charts error: %v", err))
		return
	}

	response.WriteEntity(id)
}
