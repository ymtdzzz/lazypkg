package executors

import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/ref"
)

var (
	localHashPattern  = regexp.MustCompile(`^([^\@]+)@sha256:([a-z0-9]{7})`)
	remoteHashPattern = regexp.MustCompile(`^sha256:([a-z0-9]{7})`)
)

type DockerExecutor struct {
	dc *client.Client
	rc *regclient.RegClient
}

func NewDockerExecutor() (*DockerExecutor, error) {
	dc, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	rc := regclient.New()

	return &DockerExecutor{
		dc: dc,
		rc: rc,
	}, nil
}

func (de *DockerExecutor) GetPackages(_ string) ([]*PackageInfo, error) {
	ctx := context.Background()

	var packages []*PackageInfo

	images, err := de.dc.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return packages, err
	}

	for _, image := range images {
		if len(image.RepoTags) == 0 {
			continue
		}

		// latest images only
		if !strings.Contains(image.RepoTags[0], "latest") {
			continue
		}

		// pulled images only
		if len(image.RepoDigests) == 0 {
			continue
		}

		imageName := image.RepoTags[0]
		localDigest := image.RepoDigests[0]

		r, err := ref.New(imageName)
		if err != nil {
			log.Printf("Error creating a reference for image: %s, error: %v", imageName, err)
			continue
		}
		defer de.rc.Close(ctx, r)

		m, err := de.rc.ManifestGet(ctx, r)
		if err != nil {
			log.Printf("Error getting manifest for image: %s, error: %v", imageName, err)
			continue
		}

		if pkg, err := dockerDiffPackageFromHash(imageName, localDigest, m.GetDescriptor().Digest.String()); err == nil && pkg != nil {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

func (de *DockerExecutor) Update(img, _ string, dryRun bool) error {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	errors := map[string]error{}

	wg.Add(1)
	go de.pullImage(&wg, mu, errors, img, dryRun)

	wg.Wait()

	if len(errors) > 0 {
		msg := "Some images failed to pull:"
		for img, err := range errors {
			msg = fmt.Sprintf("%s\n - %s: %v", msg, img, err)
		}
		return fmt.Errorf("%v", msg)
	}

	return nil
}

func (de *DockerExecutor) BulkUpdate(imgs []string, _ string, dryRun bool) error {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	errors := map[string]error{}

	for _, img := range imgs {
		wg.Add(1)
		go de.pullImage(&wg, mu, errors, img, dryRun)
	}

	wg.Wait()

	if len(errors) > 0 {
		msg := "Some images failed to pull:"
		for img, err := range errors {
			msg = fmt.Sprintf("%s\n - %s: %v", msg, img, err)
		}
		return fmt.Errorf("%v", msg)
	}

	return nil
}

func (de *DockerExecutor) Valid() bool {
	_, err := de.dc.Info(context.Background())
	return err == nil
}

func (de *DockerExecutor) Close() {
	if err := de.dc.Close(); err != nil {
		log.Printf("Failed to close the docker client: %v", err)
	}
}

func (de *DockerExecutor) pullImage(
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	errors map[string]error,
	img string,
	dryRun bool,
) {
	defer wg.Done()

	ctx := context.Background()

	r, err := ref.New(img)
	if err != nil {
		mu.Lock()
		errors[img] = err
		mu.Unlock()
		return
	}

	if dryRun {
		log.Printf("[dry-run] Pulling image: %s", img)
		return
	}

	out, err := de.dc.ImagePull(ctx, r.Reference, image.PullOptions{})
	if err != nil {
		mu.Lock()
		errors[img] = err
		mu.Unlock()
		return
	}
	defer out.Close()

	if _, err := io.Copy(log.Writer(), out); err != nil {
		mu.Lock()
		errors[img] = err
		mu.Unlock()
		return
	}
}

func dockerDiffPackageFromHash(imageName, localDigest, remoteDigest string) (*PackageInfo, error) {
	localMathces := localHashPattern.FindStringSubmatch(localDigest)
	if len(localMathces) < 3 {
		return nil, fmt.Errorf("invalid local digest provided: %s", localDigest)
	}
	remoteMatches := remoteHashPattern.FindStringSubmatch(remoteDigest)
	if len(remoteMatches) < 2 {
		return nil, fmt.Errorf("invalid remote digest provided: %s", remoteDigest)
	}

	ld := localMathces[2]
	rd := remoteMatches[1]

	if ld == rd {
		return nil, nil
	}

	return &PackageInfo{
		Name:       imageName,
		OldVersion: ld,
		NewVersion: rd,
	}, nil
}
