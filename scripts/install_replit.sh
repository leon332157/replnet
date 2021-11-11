#! /bin/bash
if [[ -v REPL_SLUG ]];then
    cd /home/runner/$REPL_SLUG
else
    echo "Script must be run in a repl"
    exit 1
fi

echo "Downloading replish binary to /home/runner/$REPL_SLUG"
curl -fL https://github.com/ReplDepot/replish/releases/latest/download/replish-linux-amd64 --progress-bar --compressed -o replish
if [ $? -eq 0 ]; then
    chmod +x replish
else
    echo "Download error"
    exit 1
fi
function write_exmaple_config() {
    echo "[replish]" >> .replit;
    echo 'mode = "server"' >> .replit;
}
if ! grep -q "[replish]" ".replit";then
    write_exmaple_config
else
    echo "replish field exist, skipping writing config"
fi