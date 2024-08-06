Set-Location -Path "client/account-management"  
yarn 
yarn build
$destination = "..\..\cmd\server\static\account-management"
if (-not (Test-Path -Path $destination -PathType Container)) {
    New-Item -Path $destination -ItemType Directory -Force
}
copy-item -Path ".\build\*" -Destination destination -Recurse -Force
Set-Location -Path "..\.."
 

Set-Location -Path "client/oidc-flows"  
yarn 
yarn build
$destination = "..\..\cmd\server\static\oidc-flows"
if (-not (Test-Path -Path $destination -PathType Container)) {
    New-Item -Path $destination -ItemType Directory -Force
}
copy-item -Path ".\build\*" -Destination $destination -Recurse -Force
Set-Location -Path "..\.."

 