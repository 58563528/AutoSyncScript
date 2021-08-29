#!/usr/bin/env bash

## Build 20210829-001

dir_shell=/ql/shell
. $dir_shell/share.sh

gen_pt_pin_array() {
  local envs=$(eval echo "\$JD_COOKIE")
  local array=($(echo $envs | sed 's/&/ /g'))
  local tmp1 tmp2 i pt_pin_temp
  for i in "${!array[@]}"; do
    pt_pin_temp=$(echo ${array[i]} | perl -pe "{s|.*pt_pin=([^; ]+)(?=;?).*|\1|; s|%|\\\x|g}")
    [[ $pt_pin_temp == *\\x* ]] && pt_pin[i]=$(printf $pt_pin_temp) || pt_pin[i]=$pt_pin_temp
  done
}

check_jd_cookie(){
    local test_connect="$(curl -I -s --connect-timeout 5 https://bean.m.jd.com/bean/signIndex.action -w %{http_code} | tail -n1)"
    local test_jd_cookie="$(curl -s --noproxy "*" "https://bean.m.jd.com/bean/signIndex.action" -H "cookie: $1")"
    if [ "$test_connect" -eq "302" ]; then
        [[ "$test_jd_cookie" ]] && echo "(COOKIE 有效)" || echo "(COOKIE 已失效)"
    else
        echo "(API 连接失败)"
    fi
}

dump_user_info(){
echo -e "\n## 账号用户名及 COOKIES 整理如下："
local envs=$(eval echo "\$JD_COOKIE")
local array=($(echo $envs | sed 's/&/ /g'))
    for ((m = 0; m < ${#pt_pin[*]}; m++)); do
        j=$((m + 1))
        echo -e "## 用户名 $j：${pt_pin[m]} `check_jd_cookie ${array[m]}`\nCookie$j=\"${array[m]}\""
    done
}

gen_pt_pin_array
dump_user_info