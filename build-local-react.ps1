Set-Location -Path "client/account-management"  
yarn 
yarn build
copy-item -Path ".\build\*" -Destination "..\..\cmd\server\static\account-management" -Recurse -Force
Set-Location -Path "..\.."
 

Set-Location -Path "client/keeptrack"  
yarn 
yarn build
copy-item -Path ".\build\*" -Destination "..\..\cmd\server\static\keeptrack" -Recurse -Force
Set-Location -Path "..\.."

 