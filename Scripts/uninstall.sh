#!/bin/bash

APP_NAME=borsch
BORSCH_HOME=~/borsch
BORSCH_BIN=${BORSCH_HOME}/bin
BORSCH_LIB=${BORSCH_HOME}/lib

C_BOLD_PURPLE='\033[1;35m'
C_BOLD_GREEN='\033[1;32m'
C_BOLD_BLACK='\033[1;30m'
NO_COLOR='\033[0m'

echo -e "${C_BOLD_PURPLE}==>${NO_COLOR} ${C_BOLD_BLACK}Видалення стандартної бібліотеки...${NO_COLOR}" && \
rm -rf ${BORSCH_LIB} && \
echo -e "${C_BOLD_GREEN}Готово.${NO_COLOR}" && echo && \
echo -e "${C_BOLD_PURPLE}==>${NO_COLOR} ${C_BOLD_BLACK}Видалення інтерпретатора...${NO_COLOR}" && \
rm -rf ${BORSCH_BIN} && \
echo -e "${C_BOLD_GREEN}Готово.${NO_COLOR}" && \
rm -rf ${BORSCH_HOME} && echo && \
echo -e "${C_BOLD_PURPLE}==>${NO_COLOR} ${C_BOLD_BLACK}Завершення процесу видалення.${NO_COLOR}" && echo && \
echo -e "Вилучіть змінну ${C_BOLD_BLACK}BORSCH_LIB${NO_COLOR} та обгортку ${C_BOLD_BLACK}борщ${NO_COLOR} за допомогою наступних команд:" && \
echo -e "${C_BOLD_BLACK}  unexport BORSCH_LIB${NO_COLOR}" && \
echo -e "${C_BOLD_BLACK}  unalias борщ${NO_COLOR}" && echo && \
echo 'Із профіля командної оболонки вилучіть наступні рядки, якщо вони присутні:' && \
echo -e "${C_BOLD_BLACK}  export BORSCH_LIB=${BORSCH_LIB}${NO_COLOR}" && \
echo -e "${C_BOLD_BLACK}  alias борщ="${BORSCH_BIN}/${APP_NAME}"${NO_COLOR}" && echo
