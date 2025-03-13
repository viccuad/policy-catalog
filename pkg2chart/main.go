package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/name"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

// ArtifactHubPkg represents the structure of an artifacthub-pkg.yml file
type ArtifactHubPkg struct {
	DisplayName      string            `yaml:"displayName"`
	Version          string            `yaml:"version"`
	Description      string            `yaml:"description"`
	Keywords         []string          `yaml:"keywords"`
	HomeURL          string            `yaml:"homeURL"`
	Annotations      map[string]string `yaml:"annotations"`
	ContainersImages []ContainerImage  `yaml:"containersImages"`
}

// ContainerImage represents the structure of a container image in an artifacthub-pkg.yml file
type ContainerImage struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

func main() {
	inputPath := flag.String("input", "artifacthub-pkg.yml", "Path to the artifacthub-pkg.yml file")
	outputPath := flag.String("output", "Chart.yaml", "Path to output Chart.yaml file")

	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Error: unrecognized flag or argument: %v\n", flag.Args())
		os.Exit(1)
	}

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	var pkg ArtifactHubPkg
	err = yaml.Unmarshal(data, &pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	// Delete the `kubewarden/questions-ui` annotation.
	// Rancher retrieves the questions.yaml file directly from the chart directory when listing a catalog.
	annotations := pkg.Annotations
	delete(annotations, "kubewarden/questions-ui")

	if len(pkg.ContainersImages) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no container image found in artifacthub-pkg.yml\n")
		os.Exit(1)
	}

	ociUrl := pkg.ContainersImages[0].Image
	ref, err := name.NewTag(ociUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing container image URL: %v\n", err)
		os.Exit(1)
	}

	annotations["kubewarden/registry"] = ref.Context().RegistryStr()
	annotations["kubewarden/repository"] = ref.Context().RepositoryStr()
	annotations["kubewarden/tag"] = ref.TagStr()

	metadata := chart.Metadata{
		Name:        pkg.DisplayName,
		Version:     pkg.Version,
		AppVersion:  pkg.Version,
		Description: pkg.Description,
		Keywords:    pkg.Keywords,
		Home:        pkg.HomeURL,
		Sources:     []string{pkg.ContainersImages[0].Image},
		Annotations: pkg.Annotations,
	}

	outputDir := filepath.Dir(*outputPath)
	if outputDir != "." {
		err = os.MkdirAll(outputDir, 0o755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	err = chartutil.SaveChartfile(*outputPath, &metadata)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving Chart.yaml: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", *inputPath, *outputPath)
}
