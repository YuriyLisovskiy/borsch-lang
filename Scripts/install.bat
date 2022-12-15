@echo off

setlocal
set ROOT_PACKAGE_NAME="github.com/YuriyLisovskiy/borsch-lang/Borsch"

set APP_NAME=borsch.exe
set BORSCH_HOME=%userprofile%\borsch
set BORSCH_BIN=%BORSCH_HOME%\bin
set BORSCH_LIB=%BORSCH_HOME%\lib

set C_BOLD_PURPLE=[95m
set C_BOLD_RED=[91m
set C_BOLD_GREEN=[92m
set C_BOLD_DEFAULT=[0m[1m
set NO_COLOR=[0m

:setup
  chcp 65001 1>NUL || goto done
  call :checkEnv
  call :installLib
  call :buildAndInstallInterpreter
  goto finish

:checkEnv
  echo %C_BOLD_PURPLE%==^> %C_BOLD_DEFAULT%Перевірка середовища Go...%NO_COLOR%
  where /q go || ^
  echo %C_BOLD_RED%Помилка.%NO_COLOR% && ^
  echo. && ^
  echo Не вдалося знайти систему збірки програми. && ^
  echo. && ^
  echo Див. інформацію щодо встановлення середовища Go за посиланням: && ^
  echo.   https://go.dev/doc/install && ^
  echo. && ^
  goto done
  echo %C_BOLD_GREEN%Готово.%NO_COLOR%
  echo.

:installLib
  echo %C_BOLD_PURPLE%==^> %C_BOLD_DEFAULT%Встановлення стандартної бібліотеки...%NO_COLOR%
  if not exist "%BORSCH_LIB%" mkdir %BORSCH_LIB% 1>NUL || goto done
  robocopy Lib %BORSCH_LIB% /e /nfl /ndl /njh /njs /nc /ns /np 1>NUL || goto done
  echo %C_BOLD_GREEN%Готово.%NO_COLOR%
  echo.
  echo Бібліотека міститься в каталозі %BORSCH_LIB%
  echo.

:buildAndInstallInterpreter
  setlocal
  echo %C_BOLD_PURPLE%==^>%NO_COLOR% %C_BOLD_DEFAULT%Збірка та встановлення інтерпретатора...%NO_COLOR%
  if not exist "%BORSCH_BIN%" mkdir %BORSCH_BIN% 1>NUL || endlocal && goto done
  call :getDatetime BUILD_TIME
  set LDFLAGS=-X '%ROOT_PACKAGE_NAME%/cli/build.Time=%BUILD_TIME%'
  go build -ldflags "%LDFLAGS%" -o %BORSCH_BIN%\%APP_NAME% Borsch\cli\main.go 1>NUL || endlocal && goto done
  echo %C_BOLD_GREEN%Готово.%NO_COLOR%
  echo.
  echo Інтерпретатор міститься в каталозі %BORSCH_BIN%
  echo.
  endlocal

:finish
  echo %C_BOLD_PURPLE%==^>%C_BOLD_DEFAULT% Завершення процесу встановлення.%NO_COLOR%
  call :createSymlink
  call :setEnvHome
  call :appendEnvBinToPath
  call :setEnvLib
  echo %C_BOLD_GREEN%Готово.%NO_COLOR%
  goto done

:createSymlink
  mklink %BORSCH_BIN%\борщ.exe %BORSCH_BIN%\%APP_NAME% 1>NUL || ^
  echo Запустіть команду нижче, щоб створити посилання на інтерпретатор: && ^
  echo.   mklink %BORSCH_BIN%\борщ.exe %BORSCH_BIN%\%APP_NAME% && ^
  echo.

:setEnvHome
  setx BORSCH_HOME "%BORSCH_HOME%" 1>NUL || ^
  echo Запустіть команду нижче, щоб встановити змінну середовища BORSCH_HOME: && ^
  echo.   setx BORSCH_HOME "%BORSCH_HOME%" && ^
  echo.

:appendEnvBinToPath
  setx PATH "%PATH%;%BORSCH_BIN%" 1>NUL || ^
  echo Запустіть команду нижче, щоб встановити змінну середовища BORSCH_BIN: && ^
  echo.   setx PATH "%PATH%;%BORSCH_BIN%" && ^
  echo.

:setEnvLib
  setx BORSCH_LIB "%BORSCH_LIB%" 1>NUL || ^
  echo Запустіть команду нижче, щоб встановити змінну середовища BORSCH_LIB: && ^
  echo.   setx BORSCH_LIB "%BORSCH_LIB%" && ^
  echo.

:done
  exit /b

:getDatetime
  setlocal
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
  set %~1=%month_short% %day% %year%, %current_time:~0,8%
  endlocal
