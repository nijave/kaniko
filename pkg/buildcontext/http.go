/*

Nick Venenga

A build context fetcher for HTTP

*/

package buildcontext

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GoogleContainerTools/kaniko/pkg/constants"
	"github.com/GoogleContainerTools/kaniko/pkg/util"
	"github.com/sirupsen/logrus"
)

// HTTP struct for downloading context from web url
type HTTP struct {
	context string
}

// UnpackTarFromBuildContext downloads tar and extracts it
func (h *HTTP) UnpackTarFromBuildContext() (string, error) {
	return constants.BuildContextDir, unpackTarFromHTTPSource(h.context, constants.BuildContextDir)
}

func unpackTarFromHTTPSource(url string, directory string) error {
	// Download tar file
	logrus.Debug("Downloading")
	tarPath, err := getTarFromHTTP(url, directory)
	if err != nil {
		return err
	}

	// Extract to build context
	logrus.Debug("Unpacking source context tar...")
	if err := util.UnpackCompressedTar(tarPath, directory); err != nil {
		return err
	}

	// Cleanup archive file
	logrus.Debugf("Deleting %s", tarPath)
	return os.Remove(tarPath)
}

func getTarFromHTTP(url string, directory string) (string, error) {
	filenameSplit := strings.Split(url, "/")
	filename := filenameSplit[len(filenameSplit)-1]
	httpClient := &http.Client{Timeout: time.Second * 10}

	logrus.Debugf("Downloading %s", url)
	// TODO check HTTP status code for 200 else err
	response, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	logrus.Debugf("Creating output file %s", filename)
	tarPath := filepath.Join(directory, filename)
	if err := util.CreateFile(tarPath, response.Body, 0600, 0, 0); err != nil {
		return "", err
	}

	return tarPath, nil
}
