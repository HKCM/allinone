#!/usr/bin/env bash
# github API 在未登录的情况下每小时只能调用60次

LIST=$(cat <<EOF
{
"sui":
    {"repo":"MystenLabs/sui"}
}
EOF
)

function get_info() {
    curl -s https://api.github.com/repos/$1/releases/latest | jq -cr '. | {"tag": .tag_name, "url": .html_url,"published": .published_at}'
}

function get_pub() {
    echo $1 | jq -cr .published
}

function check_date() {
    TODAY=$(date '+%Y-%m-%d')
    if [[ $1 =~ $TODAY ]];then
        send_msg $2 $3
    fi
}

function get_check_list() {
    echo $LIST | jq '. | keys' | jq -cr .[]
}

function get_repo() {
    echo $LIST | jq -cr .$1.repo
}

function send_msg() {
# send hangoutschat
    CHATROOM_URL="https://chat.googleapis.com/v1/spaces/XXXXXXX/messages?key=XXXXXXXXXXXXXX&token=XXXXXX"

    curl -s -X POST -H "Content-Type: application/json" -d '{"text":"'"$1 has new release\n$2\n<users/all>"'"}' "${CHATROOM_URL}" > /dev/null
}

get_check_list > check_list_tmp

declare -A lastest_version
while read target_repo;do
    REPO="$(get_repo ${target_repo})"
    INFO="$(get_info $REPO)"
    URL="$(echo $INFO | jq -cr .url)"
    target_version_now="$(echo $INFO | jq -cr .tag)"
    lastest_version["${target_repo}"]="$target_version_now"
    if [ -s target_last_version_info.json ]
    then
        target_last_version=$(cat ./target_last_version_info.json | jq -cr .${target_repo})
        if [ $target_version_now != ${target_last_version} ]
        then
            echo "${target_repo} has different version, last verison is ${target_last_version}, lastest version is ${target_version_now}"
            send_msg ${target_repo} $URL
        fi
    else
        PUBLISHED="$(get_pub $INFO)"
        check_date $PUBLISHED ${target_repo} $URL
    fi
done < check_list_tmp

rm check_list_tmp

json_obj="{}"

for key in "${!lastest_version[@]}"
do
    json_obj=$(echo $json_obj|jq --arg key $key --arg value "${lastest_version[$key]}" '.+{($key):$value}')
done
echo ${json_obj}

echo $json_obj|jq . > ./target_last_version_info.json
