#!/bin/bash

ROOT_PACKAGE_NAME="github.com/YuriyLisovskiy/borsch-lang/Borsch"

APP_NAME=borsch
BORSCH_HOME=~/borsch
BORSCH_BIN=${BORSCH_HOME}/bin
BORSCH_LIB=${BORSCH_HOME}/lib

C_BOLD_PURPLE='\033[1;35m'
C_BOLD_GREEN='\033[1;32m'
C_BOLD_BLACK='\033[1;30m'
NO_COLOR='\033[0m'

BUILD_TIME=$(LC_ALL=uk_UA.utf8 date '+%b %d %Y, %T')
LDFLAGS="-X '${ROOT_PACKAGE_NAME}/cli/build.Time=${BUILD_TIME}'"

echo -e "${C_BOLD_PURPLE}==>${NO_COLOR} ${C_BOLD_BLACK}Встановлення стандартної бібліотеки...${NO_COLOR}" && \
mkdir -p ${BORSCH_LIB} && \
cp -R Lib/ ${BORSCH_LIB} && \
echo -e "${C_BOLD_GREEN}Готово.${NO_COLOR}" && echo && \
echo Бібліотека міститься в каталозі ${BORSCH_LIB} && echo && \
echo -e "${C_BOLD_PURPLE}==>${NO_COLOR} ${C_BOLD_BLACK}Збірка та встановлення інтерпретатора...${NO_COLOR}" && \
mkdir -p ${BORSCH_BIN} && \
go build -ldflags "${LDFLAGS}" -o ${BORSCH_BIN}/${APP_NAME} Borsch/cli/main.go && \
echo -e "${C_BOLD_GREEN}Готово.${NO_COLOR}" && echo && \
echo Інтерпретатор міститься в каталозі ${BORSCH_BIN} && echo && \
echo -e "${C_BOLD_PURPLE}==>${NO_COLOR} ${C_BOLD_BLACK}Завершення процесу встановлення.${NO_COLOR}" && echo && \
echo У кінець профіля командної оболонки, яку ви використовуєте, додайте рядки: && \
echo -e "${C_BOLD_BLACK}  export BORSCH_LIB=${BORSCH_LIB}${NO_COLOR}" && \
echo -e "${C_BOLD_BLACK}  alias борщ="${BORSCH_BIN}/${APP_NAME}"${NO_COLOR}" && echo && \
echo Перезапустіть термінал, щоб застосувати зміни. Для цього оновіть профіль. Приклад команди для оболонки Bash: && \
echo -e "${C_BOLD_BLACK}  source ~/.bash_profile${NO_COLOR}" && echo && \
echo Для застосування змін в інших оболонках, дізнайтеся самостійно, як це зробити, або просто перезапустіть термінал. && echo
