# Management Go-App

WebAssembly application for account management, built with [go-app](https://go-app.dev/).

## Project Structure

```
pkg/go-app/management/
├── cmd/
│   ├── wasm/                    # WASM client entry point
│   ├── static_generator/        # Generates static HTML files
│   ├── fixup_static_html/       # Post-processes HTML/JS templates
│   └── server/                  # Local dev server
├── internal/
│   ├── WASMRuntime/            # WASM runtime and routing
│   ├── common/                 # Shared variables (version, build info)
│   └── services/               # Business logic services
├── web/                        # Static assets (CSS, JS, images)
└── Makefile                    # Build automation
```

## Quick Start

```bash
# Build everything and generate static files
make generate-static

# Run local dev server
make run-server

# Clean build artifacts
make clean
```

## Build Process

The build creates static files in `../../../cmd/server/static/go-app/management/static_output/`:

1. **WASM Build** - Compiles Go to WebAssembly (`app.wasm`)
2. **Static Generation** - Creates HTML/JS/CSS files
3. **Asset Copying** - Copies web assets
4. **Template Fixup** - Post-processes templates with basehref

### Output Files

```
cmd/server/static/go-app/management/static_output/
├── index.html              # Main entry point
├── index_template.html     # Template for server-side rendering
├── app.js                  # Go-app loader
├── app_template.js         # Template with basehref replacements
├── app-worker.js           # Service worker
├── app.wasm               # Compiled Go WebAssembly
└── web/
    ├── app.json           # App configuration
    ├── styles.css         # Custom styles
    └── ...                # Other assets
```

## Development Workflow

### Building Locally

```bash
cd pkg/go-app/management
make generate-static
```

### Integration with Server

After building, the fixup tool updates URLs:

```bash
# From project root
.\cmd\fixup_go_app\fixup_go_app.exe -path ./cmd/server/static/go-app/management -basehref management
```

Or use the main Makefile:

```bash
# From project root
make dev-setup
```

## Build Variables

Set via ldflags during compilation:

- `AppVersion` - Unix timestamp of build
- `BuildTime` - ISO timestamp
- `GitCommit` - Short commit hash
- `GitBranch` - Current branch name

## Customization

Override Makefile variables:

```bash
make generate-static APP_TITLE="My Dashboard" BASE_HREF="admin"
```

## Dependencies

External packages from `fluffycore-rage-identity`:

- Account management API clients
- Rage API clients
- Common services

These are managed via the main `go.mod` with replace directives.

## Integration with Mastodon Identity

The generated static files are served by the mastodon-identity server at:

- Path: `/management/`
- Files: `cmd/server/static/go-app/management/static_output/`

The server uses the template files (`*_template.*`) for server-side rendering with dynamic basehref injection.
