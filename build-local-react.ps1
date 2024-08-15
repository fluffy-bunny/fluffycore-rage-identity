Set-Location -Path "client/account-management"  
Get-ChildItem -Path "."
yarn 
yarn build
$destination = "..\..\cmd\server\static\account-management"
if (-not (Test-Path -Path $destination -PathType Container)) {
    New-Item -Path $destination -ItemType Directory -Force
}
copy-item -Path ".\account-management\*" -Destination destination -Recurse -Force
Set-Location -Path "..\.."
 

Set-Location -Path "client/oidc-flows"  
Get-ChildItem -Path "."
yarn 
yarn build
$destination = "..\..\cmd\server\static\oidc-flows"
if (-not (Test-Path -Path $destination -PathType Container)) {
    New-Item -Path $destination -ItemType Directory -Force
}
copy-item -Path ".\oidc-flows\*" -Destination $destination -Recurse -Force
Set-Location -Path "..\.."

 