#! /bin/bash
echo "Downloading replish binary to ${pwd}"
if [ -z ${REPL_SLUG} ];
then
    cd /home/runner/$REPL_SLUG
fi
curl https://github.com/leon332157/replish/releases/latest/download/replish-linux-amd64 -# --compressed -o replish
chmod +x replish

function write_exmaple_config() {
    echo "[replish]" >> .replit;
    echo 'mode = "server"' >> .replit;
}

write_exmaple_config