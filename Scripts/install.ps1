$RootPackageName = "github.com/YuriyLisovskiy/borsch-lang/Borsch"

$AppName = "borsch.exe"
$BorschHome = "$env:USERPROFILE\borsch"
$BorschBin = "$BorschHome\bin"
$BorschLib = "$BorschHome\lib"

$CBoldPurple = "[95m"
$CBoldRed = "[91m"
$CBoldGreen = "[92m"
$CBoldBlack = "[0m[1m"
$NoColor = "[0m"

Echo "$CBoldPurple==>$CBoldBlack Перевірка середовища Go...$NoColor"
$GoFound = [bool] (Get-Command -ErrorAction Ignore -Type Application go)
if (-Not $GoFound)
{
    Echo "$CBoldRed Помилка.$NoColor`n"
    Echo "Не вдалося знайти систему збірки програми.`n"
    Echo "Див. інформацію щодо встановлення середовища Go за посиланням:"
    Echo "   https://go.dev/doc/install`n"
    Exit 1
}

Echo "$CBoldGreen Готово.$NoColor`n"
Echo "$CBoldPurple==>$CBoldBlack Встановлення стандартної бібліотеки...$NoColor"
if (-Not (Test-Path -Path $BorschLib))
{
    New-Item -ItemType Directory -Path $BorschLib
}

Copy-Item -Path .\Lib -Destination $BorschLib -PassThru

Echo "$CBoldGreen Готово.$NoColor`n"
Echo "Бібліотека міститься в каталозі $BorschLib`n"

Echo "$CBoldPurple==>$CBoldBlack Збірка та встановлення інтерпретатора...$NoColor"
if (-Not (Test-Path -Path $BorschBin))
{
    New-Item -ItemType Directory -Path $BorschBin
}

$BuildTime = (Get-Date -Format "MMM dd yyyy, HH:mm:ss")
$LDFlags = "-X '$RootPackageName/cli/build.Time=$BuildTime'"
go build -ldflags "$LDFlags" -o $BorschBin\$AppName Borsch\cli\main.go

New-Item -ItemType SymbolicLink -Path $BorschBin\борщ.exe -Target $BorschBin\$AppName
[Environment]::SetEnvironmentVariable("BORSCH_LIB", "$BorschLib", "User")
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$BorschBin", "User")

Echo "$CBoldGreen Готово.$NoColor`n"
Echo "Інтерпретатор міститься в каталозі $BorschBin`n"
