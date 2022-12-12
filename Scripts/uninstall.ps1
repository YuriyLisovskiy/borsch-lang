$BorschHome = "$env:USERPROFILE\borsch"
$BorschBin = "$BorschHome\bin"
$BorschLib = "$BorschHome\lib"

$CBoldPurple = "[95m"
$CBoldGreen = "[92m"
$CBoldBlack = "[0m[1m"
$NoColor = "[0m"

Echo "$CBoldPurple==>$CBoldBlack Видалення стандартної бібліотеки...$NoColor"
Echo "Каталог, де розташована стандартна бібліотека:`n  $BorschLib"
$confirmationLib = (Read-Host "Вкажіть каталог, якщо поточний не співпадає [натисніть enter, щоб залишити поточний]:")
if ($confirmationLib -eq "")
{
    $confirmationLib = $BorschLib
}

Remove-Item $confirmationLib
Echo "$CBoldGreen Готово.$NoColor`n"

Echo "$CBoldPurple==>$NoColor $CBoldBlack Видалення інтерпретатора...$NoColor"
Echo "Каталог, де розташований інтерпретатор:`n  $BorschBin"
$confirmationBin = (Read-Host "Вкажіть каталог, якщо поточний не співпадає [натисніть enter, щоб залишити поточний]:")
if ($confirmationBin -eq "")
{
    $confirmationBin = $BorschBin
}

Remove-Item $confirmationBin
Echo "$CBoldGreen Готово.$NoColor`n"

Echo "$CBoldPurple==>$CBoldBlack Завершення процесу видалення.$NoColor`n"

(Get-Item $confirmationBin\борщ.exe).Delete()
[Environment]::SetEnvironmentVariable("BORSCH_LIB", "", "User")

Echo "$CBoldGreen Готово.$NoColor`n"
