@echo off

setlocal

set APP_NAME=borsch.exe
set BORSCH_HOME=%userprofile%\borsch
set BORSCH_BIN=%BORSCH_HOME%\bin
set BORSCH_LIB=%BORSCH_HOME%\lib

set C_BOLD_PURPLE=[95m
set C_BOLD_GREEN=[92m
set C_BOLD_DEFAULT=[0m[1m
set NO_COLOR=[0m

goto setup

:setup
  chcp 65001 > NUL
  call :removeLib
  call :removeBin
  goto finish

:removeLib
  setlocal
  echo %C_BOLD_PURPLE%==>%NO_COLOR% %C_BOLD_DEFAULT%Видалення стандартної бібліотеки...%NO_COLOR%
  echo.
  echo Каталог, де розташована стандартна бібліотека:
  echo %C_BOLD_DEFAULT%  %BORSCH_LIB%%NO_COLOR%
  echo.
  set /p borsch_lib_var=Вкажіть каталог, якщо поточний не співпадає [натисніть enter, щоб залишити поточний]: || goto done
  if "%borsch_lib_var%" == "" set borsch_lib_var=%BORSCH_LIB%
  rmdir /s /q %borsch_lib_var% || endlocal && goto done
  echo.
  echo "%C_BOLD_GREEN%Готово.%NO_COLOR%"
  echo.
  endlocal

:removeBin
  setlocal
  echo "%C_BOLD_PURPLE%==>%NO_COLOR% %C_BOLD_DEFAULT%Видалення інтерпретатора...%NO_COLOR%"
  echo.
  echo "Каталог, де розташований інтерпретатор:"
  echo "%C_BOLD_DEFAULT%  %BORSCH_BIN%%NO_COLOR%"
  echo.
  set /p borsch_bin_var=Вкажіть каталог, якщо поточний не співпадає [натисніть enter, щоб залишити поточний]: || goto done
  if "%borsch_bin_var%" == "" set borsch_bin_var=%BORSCH_BIN%
  rmdir /s /q %borsch_bin_var% || endlocal && goto done
  echo.
  echo "%C_BOLD_GREEN%Готово.%NO_COLOR%"
  echo.
  endlocal

:finish
  echo "%C_BOLD_PURPLE%==>%NO_COLOR% %C_BOLD_DEFAULT%Завершення процесу видалення.%NO_COLOR%"
  echo.
  call :removeEnvLib
  call :removeEnvBinFromPath
  call :removeEnvHome
  call :removeSymlink
  goto done

:removeEnvLib
  setx BORSCH_LIB "" 1> NUL || ^
  echo Запустіть команду нижче, щоб вилучити змінну з каталогом стандартної бібліотеки - BORSCH_LIB: && ^
  echo.   setx BORSCH_LIB "" && ^
  echo.

:removeEnvBinFromPath
  setx /m PATH "%PATH:%borsch_bin_var%;=%" 1> NUL || ^
  echo Запустіть команду нижче, щоб вилучити шлях до інтрпретатора з PATH: && ^
  echo.   setx PATH "%%PATH:%borsch_bin_var%;=%%" && ^
  echo.

:removeEnvHome:
  setx BORSCH_HOME "" 1> NUL || ^
  echo Запустіть команду нижче, щоб вилучити змінну з каталогом до кореня середовища розробки - BORSCH_HOME: && ^
  echo.   setx BORSCH_HOME "" && ^
  echo.

:removeSymlink:
  rmdir %borsch_bin_var%\борщ.exe 1> NUL || ^
  echo Запустіть команду нижче, щоб видалити посилання на інтерпретатор: && ^
  echo.   rmdir %borsch_bin_var%\борщ.exe && ^
  echo.

:done
  exit /b
