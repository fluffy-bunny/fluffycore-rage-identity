# fluffycore-rage-identity-oidc-login-go-app

fluffycore app based on https://github.com/maxence-charriere/go-app

## Developer Guide

### Prerequisites

- Go 1.25 or later
- Make

### Building the Application

#### Build All Components

To build the WebAssembly client, static generator, and server:

```bash
make all
```

This will create:

- `static_output/web/app.wasm` - The WebAssembly binary
- `static_generator` - Binary for generating static HTML files
- `server` - Binary for running the development server

#### Generate Static Output

After building, generate the static output files:

```bash
make generate-static
```

This creates all HTML files and copies web assets to the `static_output/` directory.

### Running the Development Server

Start the server:

```bash
./server
```

Or on Windows:

```powershell
.\server.exe
```

The server will start on `http://localhost:3557`

### Viewing the Application

Open your browser and navigate to:

```
http://localhost:3557/oidc-login/
```

**Note:** The trailing slash is important!

### Customizing the Base Path

The application is configured to run at `/oidc-login/` by default. To host it at a different path:

#### Option 1: Rebuild with Custom Path

Edit the `Makefile` and change the `BASE_HREF` variable:

```makefile
BASE_HREF := my-custom-path
```

Then rebuild:

```bash
make all
make generate-static
```

#### Option 2: Modify Generated Files

If you want to change the path after generation, you need to update these files in `static_output/`:

1. **app.js** - Find and replace `/oidc-login` with your path (e.g., `/my-root` or `/my/nested-root`)
2. **app-worker.js** - Find and replace `/oidc-login` with your path
3. **manifest.webmanifest** - Find and replace `/oidc-login` with your path
4. **All HTML files** - Update the `<base href="/oidc-login/">` tag to your new path

**Example:** To change from `/oidc-login` to `/my-app`:

```bash
# In static_output/
sed -i 's|/oidc-login|/my-app|g' app.js app-worker.js manifest.webmanifest *.html
```

On Windows PowerShell:

```powershell
# In static_output/
Get-ChildItem -Include app.js,app-worker.js,manifest.webmanifest,*.html | ForEach-Object {
    (Get-Content $_) -replace '/oidc-login', '/my-app' | Set-Content $_
}
```

### Build Versioning

Each build includes version information embedded via ldflags:

- **AppVersion**: Unix timestamp for cache busting
- **BuildTime**: Human-readable build timestamp
- **GitCommit**: Short commit SHA
- **GitBranch**: Current branch name

These are displayed in the UI and used for cache control.

### Project Structure

```
├── cmd/
│   ├── server/          # Development server
│   ├── static_generator/ # Static HTML generator
│   └── wasm/            # WebAssembly main package
├── internal/
│   ├── contracts/       # Service interfaces
│   ├── services/        # Service implementations
│   │   ├── App/         # Main app service
│   │   └── composers/   # Page components
│   └── WASMRuntime/     # WASM runtime setup
├── web/                 # Static assets (CSS, JS, images)
├── static_output/       # Generated distribution files
└── Makefile            # Build configuration
```

### Development Tips

- Use `make clean` to remove all build artifacts
- The `static_output/` directory is the complete distributable package
- Modify `web/styles.css` for styling changes
- Add new pages in `internal/services/composers/`
- Update routes in `internal/WASMRuntime/services.go`

## Rebasing

### First rebase the static files

```shell
.\rebase_static.exe --input=./static_output --output=./rebased_output --old-base=/oidc-login --new-base=/my/nested/app

./rebase_static -input=./static_output -output=./rebased_output -old-base=/oidc-login -new-base=/my/nested/app
```

### Then run server with custom settings

```shell
./server.exe --dir=./rebased_output --base=/my/nested/app -port=3558

./server -dir=./rebased_output -base=/my/nested/app -port=3558

```
