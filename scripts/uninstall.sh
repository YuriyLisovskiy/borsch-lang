#!/bin/bash

APP_NAME=borsch
BORSCH_HOME=~/borsch
BORSCH_BIN=${BORSCH_HOME}/bin
BORSCH_LIB=${BORSCH_HOME}/lib

C_BOLD_PURPLE='\033[1;35m'
C_BOLD_GREEN='\033[1;32m'
C_BOLD_DEFAULT='\033[1m'
NO_COLOR='\033[0m'

echo -e "${C_BOLD_PURPLE}==>${NO_COLOR} ${C_BOLD_DEFAULT}Видалення стандартної бібліотеки...${NO_COLOR}" && echo && \
echo "Каталог, де розташована стандартна бібліотека:" && \
echo -e "${C_BOLD_DEFAULT}  ${BORSCH_LIB}${NO_COLOR}" && echo && \
read -p "Вкажіть каталог, якщо поточний не співпадає [натисніть enter, щоб залишити поточний]: " borsch_lib_var
if [ -z "$borsch_lib_var" ]
then
  borsch_lib_var=${BORSCH_LIB}
fi

rm -rf ${borsch_lib_var} && \
echo && echo -e "${C_BOLD_GREEN}Готово.${NO_COLOR}" && echo && \

echo -e "${C_BOLD_PURPLE}==>${NO_COLOR} ${C_BOLD_DEFAULT}Видалення інтерпретатора...${NO_COLOR}" && echo && \
echo "Каталог, де розташований інтерпретатор:" && \
echo -e "${C_BOLD_DEFAULT}  ${BORSCH_BIN}${NO_COLOR}" && echo && \
read -p "Вкажіть каталог, якщо поточний не співпадає [натисніть enter, щоб залишити поточний]: " borsch_bin_var
if [ -z "$borsch_bin_var" ]
then
  borsch_bin_var=${BORSCH_BIN}
fi

rm -rf ${borsch_bin_var} && \
echo && echo -e "${C_BOLD_GREEN}Готово.${NO_COLOR}" && echo && \

echo -e "${C_BOLD_PURPLE}==>${NO_COLOR} ${C_BOLD_DEFAULT}Завершення процесу видалення.${NO_COLOR}" && echo && \
echo -e "Вилучіть змінну ${C_BOLD_DEFAULT}BORSCH_LIB${NO_COLOR} та обгортку ${C_BOLD_DEFAULT}борщ${NO_COLOR} за допомогою наступних команд:" && \
echo -e "${C_BOLD_DEFAULT}  unexport BORSCH_LIB${NO_COLOR}" && \
echo -e "${C_BOLD_DEFAULT}  unalias борщ${NO_COLOR}" && echo && \
echo 'Із профіля командної оболонки вилучіть наступні рядки, якщо вони присутні:' && \
echo -e "${C_BOLD_DEFAULT}  export BORSCH_LIB=${borsch_lib_var}${NO_COLOR}" && \
echo -e "${C_BOLD_DEFAULT}  alias борщ="${borsch_bin_var}/${APP_NAME}"${NO_COLOR}" && echo
