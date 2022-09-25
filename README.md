# yamlctl

An experimental tool to modify YAMLs without losing (most of) comment lines.

```bash
yamlctl edit --set $.foo.bar=baz <input.yaml >output.yaml
```

Inline comments and some whitespaces are not preserved.

## Install
```bash
go install github.com/AkihiroSuda/yamlctl/cmd/yamlctl@master
```

## Usage
### `yamlctl editable`
```bash
yamlctl editable input.yaml
```

Prints `true` if the YAML is well formatted for `yamlctl edit`.

Misformatted YAMLs are not safely editable, but can be edited forcibly with `yamlctl edit --force`.

### `yamlctl edit`
```bash
yamlctl edit --set $.foo.bar=baz <input.yaml >output.yaml
```

Set `-w` to write back the result to `input.yaml` directly:
```bash
yamlctl edit -w --set $.foo.bar=baz input.yaml
```

Set `--bak` to create backup files like `input.yaml.bak.0`, `input.yaml.bak.1`, ...

Set `--force` to edit a misformatted YAML forcibly.

### `yamlctl query`
```bash
yamlctl query $.foo.bar input.yaml
```

### `yamlctl yaml2json`
```bash
yamlctl yaml2json input.yaml
```

## Go library (`yamlutil`)
https://pkg.go.dev/github.com/AkihiroSuda/yamlctl/pkg/yamlutil

See the unit tests for the usage.

## Similar projects
- https://github.com/vmware-archive/go-yaml-edit (abandoned)
