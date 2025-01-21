在 PowerShell 中支持 
```shell
cd autocomplete/powershell/
mkdir autocomplete
cwgo completion powershell | Out-File autocomplete/cwgo.ps1
& autocomplete/cwgo.ps1
```
