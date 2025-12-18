# FluffyCore Rage Identity - Development Environment Setup Script
# This script installs all required tools for Windows 11 development

param(
    [switch]$SkipOptional = $false,
    [switch]$Force = $false,
    [switch]$Verbose = $false
)

# Set error action preference
$ErrorActionPreference = "Stop"

# Colors for output
$Green = "Green"
$Yellow = "Yellow" 
$Red = "Red"
$Cyan = "Cyan"

function Write-Header {
    param([string]$Message)
    Write-Host "`n=== $Message ===" -ForegroundColor $Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "✅ $Message" -ForegroundColor $Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "⚠️  $Message" -ForegroundColor $Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "❌ $Message" -ForegroundColor $Red
}

function Write-Info {
    param([string]$Message)
    Write-Host "ℹ️  $Message" -ForegroundColor $Cyan
}

function Test-AdminPrivileges {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Test-CommandExists {
    param([string]$Command)
    try {
        Get-Command $Command -ErrorAction Stop | Out-Null
        return $true
    }
    catch {
        return $false
    }
}

function Install-Chocolatey {
    Write-Header "Installing Chocolatey Package Manager"
    
    if (Test-CommandExists "choco") {
        Write-Success "Chocolatey is already installed"
        return
    }
    
    try {
        Write-Info "Installing Chocolatey..."
        Set-ExecutionPolicy Bypass -Scope Process -Force
        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
        Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
        
        # Refresh environment variables
        $env:PATH = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
        
        Write-Success "Chocolatey installed successfully"
    }
    catch {
        Write-Error "Failed to install Chocolatey: $($_.Exception.Message)"
        throw
    }
}

function Install-Go {
    Write-Header "Installing Go Programming Language"
    
    if (Test-CommandExists "go") {
        $version = (go version 2>$null) -replace "go version go", "" -replace " .*", ""
        if ($version -and [Version]$version -ge [Version]"1.25") {
            Write-Success "Go $version is already installed and meets requirements"
            return
        }
        else {
            Write-Warning "Go $version is installed but version 1.25+ is required"
        }
    }
    
    try {
        Write-Info "Installing Go via Chocolatey..."
        choco install golang -y --force
        
        # Refresh environment variables
        $env:PATH = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
        
        Write-Success "Go installed successfully"
    }
    catch {
        Write-Error "Failed to install Go: $($_.Exception.Message)"
        throw
    }
}

function Install-GitForWindows {
    Write-Header "Installing Git for Windows (includes Git Bash and Unix tools)"
    
    if (Test-CommandExists "git") {
        Write-Success "Git is already installed"
        return
    }
    
    try {
        Write-Info "Installing Git for Windows via Chocolatey..."
        choco install git -y --force
        
        # Refresh environment variables
        $env:PATH = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
        
        Write-Success "Git for Windows installed successfully"
    }
    catch {
        Write-Error "Failed to install Git for Windows: $($_.Exception.Message)"
        throw
    }
}

function Install-DockerDesktop {
    Write-Header "Installing Docker Desktop"
    
    if (Test-CommandExists "docker") {
        Write-Success "Docker is already installed"
        return
    }
    
    try {
        Write-Info "Installing Docker Desktop via Chocolatey..."
        choco install docker-desktop -y --force
        
        Write-Warning "Docker Desktop requires a system restart to complete installation"
        Write-Info "Please restart your computer after this script completes"
        
        Write-Success "Docker Desktop installed successfully"
    }
    catch {
        Write-Error "Failed to install Docker Desktop: $($_.Exception.Message)"
        Write-Warning "You may need to install Docker Desktop manually from docker.com"
    }
}

function Install-Make {
    Write-Header "Installing Make for Windows"
    
    if (Test-CommandExists "make") {
        Write-Success "Make is already installed"
        return
    }
    
    try {
        Write-Info "Installing Make via Chocolatey..."
        choco install make -y --force
        
        # Refresh environment variables
        $env:PATH = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
        
        Write-Success "Make installed successfully"
    }
    catch {
        Write-Error "Failed to install Make: $($_.Exception.Message)"
        throw
    }
}

function Install-Protoc {
    Write-Header "Installing Protocol Buffers Compiler (protoc)"
    
    if (Test-CommandExists "protoc") {
        Write-Success "protoc is already installed"
        return
    }
    
    try {
        Write-Info "Installing protoc via Chocolatey..."
        choco install protoc -y --force
        
        # Refresh environment variables
        $env:PATH = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
        
        Write-Success "protoc installed successfully"
    }
    catch {
        Write-Error "Failed to install protoc: $($_.Exception.Message)"
        throw
    }
}

function Install-GoProtobufTools {
    Write-Header "Installing Go Protocol Buffer Tools"
    
    $tools = @(
        "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest",
        "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest", 
        "google.golang.org/protobuf/cmd/protoc-gen-go@latest",
        "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
        "github.com/fluffy-bunny/fluffycore/protoc-gen-go-fluffycore-di/cmd/protoc-gen-go-fluffycore-di@latest",
        "github.com/fluffy-bunny/fluffycore/protoc-gen-go-fluffycore-di/cmd/protoc-gen-go-fluffycore-nats@latest"
    )
    
    foreach ($tool in $tools) {
        try {
            Write-Info "Installing $tool..."
            go install $tool
        }
        catch {
            Write-Warning "Failed to install $tool - will retry later"
        }
    }
    
    Write-Success "Go protobuf tools installation completed"
}

function Install-OptionalTools {
    if ($SkipOptional) {
        Write-Info "Skipping optional tools installation"
        return
    }
    
    Write-Header "Installing Optional Development Tools"
    
    # mkcert for local SSL certificates
    try {
        if (-not (Test-CommandExists "mkcert")) {
            Write-Info "Installing mkcert..."
            choco install mkcert -y --force
        } else {
            Write-Success "mkcert is already installed"
        }
    }
    catch {
        Write-Warning "Failed to install mkcert (optional)"
    }
    
    # ngrok for HTTPS tunneling
    try {
        if (-not (Test-CommandExists "ngrok")) {
            Write-Info "Installing ngrok..."
            choco install ngrok -y --force
        } else {
            Write-Success "ngrok is already installed"
        }
    }
    catch {
        Write-Warning "Failed to install ngrok (optional)"
    }
}

function Test-Installation {
    Write-Header "Verifying Installation"
    
    $tests = @(
        @{ Command = "go"; Name = "Go"; Required = $true },
        @{ Command = "git"; Name = "Git"; Required = $true },
        @{ Command = "docker"; Name = "Docker"; Required = $true },
        @{ Command = "make"; Name = "Make"; Required = $true },
        @{ Command = "protoc"; Name = "protoc"; Required = $true },
        @{ Command = "mkcert"; Name = "mkcert"; Required = $false },
        @{ Command = "ngrok"; Name = "ngrok"; Required = $false }
    )
    
    $allGood = $true
    
    foreach ($test in $tests) {
        if (Test-CommandExists $test.Command) {
            try {
                $version = & $test.Command --version 2>$null | Select-Object -First 1
                Write-Success "$($test.Name): $version"
            }
            catch {
                Write-Success "$($test.Name): Available"
            }
        }
        else {
            if ($test.Required) {
                Write-Error "$($test.Name): Not found (REQUIRED)"
                $allGood = $false
            }
            else {
                Write-Warning "$($test.Name): Not found (optional)"
            }
        }
    }
    
    return $allGood
}

function Show-PostInstallInstructions {
    Write-Header "Post-Installation Instructions"
    
    Write-Info "1. Restart your computer if Docker Desktop was installed"
    Write-Info "2. Open a new PowerShell or Git Bash terminal to refresh environment variables"
    Write-Info "3. Run makefiles from Git Bash for best compatibility with Unix commands"
    Write-Info "4. Verify your setup by running the verification commands in the README"
    Write-Info "5. Configure your .env.secrets file with your OAuth credentials"
    
    Write-Success "Development environment setup completed!"
    Write-Info "You can now run: make help"
}

# Main execution
try {
    Write-Header "FluffyCore Rage Identity Development Environment Setup"
    Write-Info "This script will install all required development tools for Windows 11"
    
    if (-not (Test-AdminPrivileges)) {
        Write-Warning "This script requires administrator privileges for some installations"
        Write-Info "Please run PowerShell as Administrator and try again"
        exit 1
    }
    
    # Install tools in order
    Install-Chocolatey
    Install-Go
    Install-GitForWindows  
    Install-DockerDesktop
    Install-Make
    Install-Protoc
    Install-GoProtobufTools
    Install-OptionalTools
    
    # Test everything
    $success = Test-Installation
    
    if ($success) {
        Show-PostInstallInstructions
    }
    else {
        Write-Error "Some required tools failed to install. Please check the errors above and install manually."
        exit 1
    }
}
catch {
    Write-Error "Setup failed: $($_.Exception.Message)"
    Write-Info "Please check the error above and try running individual installation commands manually"
    exit 1
}