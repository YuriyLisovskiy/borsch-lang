@echo off

set ROOT_PACKAGE_NAME="github.com/YuriyLisovskiy/borsch-lang/Borsch"

set APP_NAME=borsch.exe
set BORSCH_HOME=%userprofile%\borsch
set BORSCH_BIN=%BORSCH_HOME%\bin
set BORSCH_LIB=%BORSCH_HOME%\lib

set C_BOLD_PURPLE=[95m
set C_BOLD_RED=[91m
set C_BOLD_GREEN=[92m
set C_BOLD_BLACK=[0m[1m
set NO_COLOR=[0m

for /f "delims=" %%a in ('wmic OS Get localdatetime ^| find "."') do set dt=%%a
set year=%dt:~0,4%
set month=%dt:~4,2%
set day=%dt:~6,2%
if %month%==01 set month_short=Jan
if %month%==02 set month_short=Feb
if %month%==03 set month_short=Mar
if %month%==04 set month_short=Apr
if %month%==05 set month_short=May
if %month%==06 set month_short=Jun
if %month%==07 set month_short=Jul
if %month%==08 set month_short=Aug
if %month%==09 set month_short=Sep
if %month%==10 set month_short=Oct
if %month%==11 set month_short=Nov
if %month%==12 set month_short=Dec
for /f "tokens=1-3 delims=/:" %%a in ("%TIME%") do (set current_time=%%a:%%b:%%c)

set BUILD_TIME=%month_short% %day% %year%, %current_time:~0,8%
set LDFLAGS=-X '%ROOT_PACKAGE_NAME%/cli/build.Time=%BUILD_TIME%'

chcp 65001> NUL

echo %C_BOLD_PURPLE%==^> %C_BOLD_BLACK%Перевірка середовища Go...%NO_COLOR%
where /q go || ^
echo %C_BOLD_RED%Помилка.%NO_COLOR% && ^
echo. && ^
echo Не вдалося знайти систему збірки програми. && ^
echo. && ^
echo Див. інформацію щодо встановлення середовища Go за посиланням: && ^
echo.   https://go.dev/doc/install && ^
echo. && ^
EXIT /B
echo %C_BOLD_GREEN%Готово.%NO_COLOR%

echo %C_BOLD_PURPLE%==^> %C_BOLD_BLACK%Встановлення стандартної бібліотеки...%NO_COLOR%
if not exist "%BORSCH_LIB%" mkdir %BORSCH_LIB%
robocopy Lib %BORSCH_LIB% /E /NFL /NDL /NJH /NJS /nc /ns /np
echo %C_BOLD_GREEN%Готово.%NO_COLOR%
echo.
echo Бібліотека міститься в каталозі %BORSCH_LIB%
echo.

echo %C_BOLD_PURPLE%==^>%NO_COLOR% %C_BOLD_BLACK%Збірка та встановлення інтерпретатора...%NO_COLOR%
if not exist "%BORSCH_BIN%" mkdir %BORSCH_BIN%
go build -ldflags "%LDFLAGS%" -o %BORSCH_BIN%\%APP_NAME% Borsch\cli\main.go
echo @C:\Users\YuriyLisovskiy\borsch\bin\borsch.exe  > C:\Windows\System32\борщ.bat
setx /M PATH "%PATH%;%BORSCH_BIN%" > NUL
echo %C_BOLD_GREEN%Готово.%NO_COLOR%
echo.
echo Інтерпретатор міститься в каталозі %BORSCH_BIN%
echo.
echo Перезапустіть термінал, щоб застосувати зміни.
echo.