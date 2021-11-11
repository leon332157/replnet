#! /bin/bash
if [[ -v ${REPL_SLUG} ]];then
    cd /home/runner/$REPL_SLUG
else
    echo "Script must be run in a repl"
    exit 1
fi

echo "Downloading replish binary to /home/runner/$REPL_SLUG"
curl -fL https://github.com/ReplDepot/replish/releases/latest/download/replish-linux-amd64 -# --compressed -o replish
chmod +x replish

function write_exmaple_config() {
    echo "[replish]" >> .replit;
    echo 'mode = "server"' >> .replit;
}

write_exmaple_config