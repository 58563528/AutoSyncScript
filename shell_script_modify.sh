function jddj(){
    # 备份cookie文件
    [[ -f /dailycheckin/jddj/jddj_cookie.js ]] && cp -rf /dailycheckin/jddj/jddj_cookie.js /dailycheckin/backup_jddj_cookie.js
    # clone
    rm -rf /dailycheckin/jddj && git clone https://ghproxy.com/https://github.com/passerby-b/JDDJ /dailycheckin/jddj
    # 下载自定义cookie文件地址,如私密的gist地址,需修改
    jddj_cookiefile="https://ghproxy.com/https://raw.githubusercontent.com/passerby-b/JDDJ/main/jddj_cookie.js"
    curl -so /dailycheckin/jddj/jddj_cookie.js $jddj_cookiefile
    # 下载cookie文件失败时从备份恢复
    test $? -eq 0 || cp -rf /dailycheckin/jddj/backup_jddj_cookie.js /dailycheckin/backup_jddj_cookie.js
    # 获取js文件中cron字段设置定时任务
    for jsname in $(ls /dailycheckin/jddj | grep -E "js$" | tr "\n" " "); do
        jsname_cn="$(grep "cron" /dailycheckin/jddj/$jsname | grep -oE "/?/?tag\=.*" | cut -d"=" -f2)"
        jsname_log="$(echo /dailycheckin/jddj/$jsname | sed 's;^.*/\(.*\)\.js;\1;g')"
        jsnamecron="$(cat /dailycheckin/jddj/$jsname | grep -oE "/?/?cron \".*\"" | cut -d\" -f2)"
        test -z "$jsname_cn" && jsname_cn=$jsname_log
        test -z "$jsnamecron" || echo "# $jsname_cn" >> /dailycheckin/docker/merged_list_file.sh
        test -z "$jsnamecron" || echo "$jsnamecron node /dailycheckin/jddj/$jsname >> /dailycheckin/logs/$jsname_log.log 2>&1" >> /dailycheckin/docker/merged_list_file.sh
    done
}

function didi(){
    # clone
    rm -rf /dailycheckin/didi && git clone https://ghproxy.com/https://github.com/passerby-b/didi_fruit /dailycheckin/didi
    curl -so /dailycheckin/didi/sendNotify.js https://ghproxy.com/https://raw.githubusercontent.com/passerby-b/JDDJ/main/sendNotify.js
    # 获取js文件中cron字段设置定时任务
    echo "10 0,8,12,18 * * * node /dailycheckin/didi/dd_fruit.js >> /dailycheckin/logs/dd_fruit.log 2>&1" >> /dailycheckin/docker/merged_list_file.sh
}

function main(){
    jddj
    didi
}

main
