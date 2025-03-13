# pkg2chart

A simple utility to convert ArtifactHub package descriptions (`artifacthub-pkg.yml`) to Helm Chart metadata files (`Chart.yaml`) containing information needed by the Rancher UI
to catalog policies.

Note: this tool does not generate a full Helm Chart, only the metadata file.

## Usage

```bash
pkg2chart [--input artifacthub-pkg.yml] [--output Chart.yaml]
```

### Options

- `--input` - Path to the ArtifactHub package description file (default: `artifacthub-pkg.yml`)
- `--output` - Path to the output Helm Chart metadata file (default: `Chart.yaml`)
