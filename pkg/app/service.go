package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	windowsPlatform = "windows"
	apkPlatform     = "apk"

	ReleaseTypeRelease = "release"

	GithubOwner = "BacoFoods"
)

type service struct {
	gitToken      string
	gitRepository string
}

func NewService(gitToken, gitRepository string) Service {
	return &service{
		gitToken:      gitToken,
		gitRepository: gitRepository,
	}
}

// DownloadWindows downloads the Windows executable from the repository.
func (s *service) DownloadWindows(version string) (string, io.ReadCloser, int64, error) {
	return s.download(windowsPlatform, version)
}

// DownloadAPK downloads the APK from the repository.
func (s *service) DownloadAPK(version string) (string, io.ReadCloser, int64, error) {
	return s.download(apkPlatform, version)
}

// download downloads the asset from the repository.
// It returns the filename and file data.
func (s *service) download(platform, version string) (string, io.ReadCloser, int64, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var release *github.RepositoryRelease
	var err error

	logrus.Infof("Getting release %s from %s", version, s.gitRepository)

	// Get the latest release
	if version == "latest" {
		release, err = s.getLatestRelease(client)
	} else {
		release, err = s.findRelease(client, version)
	}

	if err != nil {
		return "", nil, 0, err
	}

	if release == nil {
		return "", nil, 0, fmt.Errorf("release not found")
	}

	logrus.Infof("Found release %s", *release.Name)

	// Get the latest release's assets
	assets, _, err := client.Repositories.ListReleaseAssets(ctx, GithubOwner, s.gitRepository, release.GetID(), nil)
	if err != nil {
		return "", nil, 0, err
	}

	// Search for the asset by name
	for _, asset := range assets {
		if assetMatches(asset, platform, ReleaseTypeRelease) {
			// Download the asset into a buffer
			data, err := s.downloadAsset(client, *asset.ID)
			if err != nil {
				return "", nil, 0, err
			}

			return asset.GetName(), data, int64(asset.GetSize()), nil
		}
	}

	return "", nil, 0, fmt.Errorf("asset not found")
}

// getLatestRelease gets the latest release from the repository.
func (s *service) getLatestRelease(client *github.Client) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	release, _, err := client.Repositories.GetLatestRelease(ctx, GithubOwner, s.gitRepository)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func assetMatches(asset *github.ReleaseAsset, platform, releaseType string) bool {
	name := strings.ToLower(*asset.Name)

	return strings.Contains(name, platform) && strings.Contains(name, releaseType)
}

// findRelease finds the release by name.
func (s *service) findRelease(client *github.Client, version string) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	releases, _, err := client.Repositories.ListReleases(ctx, GithubOwner, s.gitRepository, nil)
	if err != nil {
		return nil, err
	}

	// Search for the release by name
	for _, r := range releases {
		if r.GetName() == version {
			return r, err
		}
	}

	return nil, err
}

// downloadAsset downloads the asset from the provided URL and returns the file data.
func (s *service) downloadAsset(client *github.Client, id int64) (io.ReadCloser, error) {
	ctx := context.Background()
	resp, _, err := client.Repositories.DownloadReleaseAsset(ctx, GithubOwner, s.gitRepository, id, http.DefaultClient)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
