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


ls -oTFGtcr,as

echo -e "that was an${RED} ls -oTFGtcr,as ${RESET}:: -c means the time when the file status was last changed -c"
echo -e "${GREEN}do a naked -t if you instead want it sorted by time of last modification of contents ${RED}only${RESET}"
echo "... maybe also try fl and gr"
