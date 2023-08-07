- based on https://github.com/CheckPointSW/terraform-provider-checkpoint/tree/master/commands
- added login and save session ID when no sid.json found

Ad-hoc usage:
- e.g. cd ./commands/publish; go run publish.go

Windows build (powershell):
```powershell
# install go
winget install GoLang.Go
# clone repo
cd $env:TEMP; mkdir build-cp-cli; cd build-cp-cli; git clone https://github.com/mkol5222/my-cp-mgmt-commands; cd my-cp-mgmt-commands
# build commands
@("publish","discard", "logout", "login") | %{ Write-Host "Building $_"; pushd commands/$_; go get; go build "$_.go"; popd }
# show commands
gci -Rec *.exe
```