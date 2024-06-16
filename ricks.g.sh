#!/bin/bash

# Define color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
RESET='\033[0m'
# echo -e "${RED}This text is red${RESET}"
# echo -e "${GREEN}This text is green${RESET}"

gls -oaFGpstr --color=auto --time=status --time-style=long-iso --group-directories-first

echo -e "that was a${RED} gls -oaFGpstr --color=auto --time=${CYAN}status${RED} --time-style=long-iso --group-directories-first"

echo -e "${RESET}switch:${RED} --group-directories-first${RESET}; can be augmented with a ${RED}--sort option${RESET}"

echo -e "--time=WORD (WORD = atime, access, use; ctime, ${CYAN}status${RESET}; mtime, modification; birth, creation${RESET})"

echo -e "edit and change ${CYAN}status${RESET} to ${RED}mtime${RESET} if you want it sorted by content modification ${RED}only${RESET}"

echo "... maybe also try fl and gr"
