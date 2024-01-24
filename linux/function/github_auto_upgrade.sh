#!/usr/bin/env bash
set -e
# author: karl
# date; 2023-11-24
# 关键词: github自动升级

repo_owner="MystenLabs"
repo_name="sui"
download_dir="./"
download_url=""

function get_latest_info() {
    release_info=$(curl -s "https://api.github.com/repos/$repo_owner/$repo_name/releases/latest")

    download_url=$(echo "$release_info" | grep -o '"browser_download_url": ".*sui-mainnet-.*ubuntu-x86_64.*\.tgz"' | sed 's/"browser_download_url": "//;s/"$//' | head -n 1)

    echo "Latest Download URL: ${download_url}"
}

function download_pkg(){
    pkg_name=$(echo "$download_url" | awk -F'/' '{print $NF}')
    if [ -f ${pkg_name} ]; then
        echo "${pkg_name} already downloaded..."
        return
    fi

    if [ -n "$download_url" ]; then
        echo "Start download: $download_url"

        curl -L -O "$download_url"
        
        echo "Downloaded, Package is ${pkg_name}"
    else
        echo "Can not find 'ubuntu-x86_64' URL"
        exit 1
    fi

}

function confirm() {
    local prompt_message=$1
    [ -n prompt_message ] || prompt_message="Are You Sure?"
    prompt_message="${prompt_message} [YES/n\]:"
    echo ${prompt_message} 
    read -r input

    case $input in
        [yY][eE][sS]|[yY])
            # echo "Yes"
            # Other function
            echo "Confirmed.."
            ;;

        [nN][oO]|[nN])
            echo "Canceled..."
            exit 0
            ;;
        *)
            echo "Invalid input..."
            exit 1
            ;;
    esac
}

function operation(){
    # Stop the appication
    echo "Stopping application.."
    cd /data/work/sui-fullnode && bash app.sh stop
    sleep 10
    # Backup the old app
    echo "Start backup old version..."
    cd /data/work && mv target target.bak.$(date "+%Y-%m-%d-%H-%M-%S")
    # 
    echo "Start unzip new version..."
    tar -xzvf ${pkg_name}
    # Start the app
    ln -sf /data/work/target/release/sui-node-ubuntu-x86_64 /data/work/sui-fullnode/sui-node
    sleep 5
    echo "Start new version app..."
    cd /data/work/sui-fullnode && bash app.sh start
}

# TODO
#check_md5() {}

function main(){
    get_latest_info
    confirm "Do you want to download latest version?"
    download_pkg
    confirm "Do you want to relase new version?"
    operation
}

main