# Tasks: GitHub Actions + GoReleaser Workflow

## Estado: ✅ COMPLETADO

### Tareas implementadas

- [x] 1. Crear `.goreleaser.yaml` en la raíz del proyecto
- [x] 2. Crear `.github/workflows/build.yml`
- [x] 3. Verificar compilación local con `go build -o mcp-trello .`

### Detalles de implementación

#### 1. `.goreleaser.yaml`
- Binary: `mcp-trello`
- Main: `.`
- Platforms: linux (amd64, arm64), darwin (amd64, arm64), windows (amd64)
- Formats: tar.gz (Unix), zip (Windows)
- Checksums: SHA256
- Ldflags: version, commit, date

#### 2. `.github/workflows/build.yml`
- Triggers: push a main, pull_request a main, workflow_dispatch, push de tags v*
- Jobs: test (go test ./...), build
- Usa goreleaser/goreleaser-action@v6
- Solo hace release cuando es tag (startsWith(github.ref, 'refs/tags/v'))

#### 3. Verificación de compilación
- ✅ `go build -o mcp-trello .` ejecutado correctamente
