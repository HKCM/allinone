#!/usr/bin/env bash
function repeat() { 
    local i=1
    local count=$1
    shift;
    while [ $i -le ${count} ]
    do
        $@ && return;
        let i++;
        sleep 30;
    done 
}

# repeat 5 echo "Hello World!" # 指定重试次数为5