#!/bin/bash
# Author: Jrohy
# github: https://github.com/Jrohy/trojan

#Definir variables de operación, 0 es no, 1 es sí
HELP=0

REMOVE=0

UPDATE=0

#DOWNLAOD_URL="https://github.com/ThonyDroidYT/trojan/releases/download/"
DOWNLAOD_URL="https://github.com/Jrohy/trojan/releases/download/"

#VERSION_CHECK="https://api.github.com/repos/ThonyDroidYT/trojan/releases/latest"
VERSION_CHECK="https://api.github.com/repos/Jrohy/trojan/releases/latest"

#SERVICE_URL="https://raw.githubusercontent.com/ThonyDroidYT/trojan/master/asset/trojan-web.service"
SERVICE_URL="https://raw.githubusercontent.com/Jrohy/trojan/master/asset/trojan-web.service"

[[ -e /var/lib/trojan-manager ]] && UPDATE=1

#Centos Cancelar temporalmente el alias
[[ -f /etc/redhat-release && -z $(echo $SHELL|grep zsh) ]] && unalias -a

[[ -z $(echo $SHELL|grep zsh) ]] && SHELL_WAY="bash" || SHELL_WAY="zsh"

#######color code########
RED="31m"
GREEN="32m"
YELLOW="33m"
BLUE="36m"
FUCHSIA="35m"

colorEcho(){
    id=$(cat /etc/newadm/idioma)
    [ -z "$id" ] && id=es
    COLOR=$1
    #echo -e "\033[${COLOR}${@:2}\033[0m"
    echo -e "\033[${COLOR}$(source trans -e bing -b zh:${id} "${@:2}")\033[0m"
}

#######get params#########
while [[ $# > 0 ]];do
    KEY="$1"
    case $KEY in
        --remove)
        REMOVE=1
        ;;
        -h|--help)
        HELP=1
        ;;
        *)
                # unknown option
        ;;
    esac
    shift # past argument or value
done
#############################

help(){
    echo "bash $0 [-h|--help] [--remove]"
    echo "  -h, --help           Mostrar ayuda"
    echo "      --remove         Remover trojan"
    return 0
}

removeTrojan() {
    #Eliminar trojan
    rm -rf /usr/bin/trojan >/dev/null 2>&1
    rm -rf /usr/local/etc/trojan >/dev/null 2>&1
    rm -f /etc/systemd/system/trojan.service >/dev/null 2>&1

    #Eliminar el programa de gestión de trojan
    rm -f /usr/local/bin/trojan >/dev/null 2>&1
    rm -rf /var/lib/trojan-manager >/dev/null 2>&1
    rm -f /etc/systemd/system/trojan-web.service >/dev/null 2>&1

    systemctl daemon-reload

    #Eliminar la base de datos dedicada a trojan
    docker rm -f trojan-mysql trojan-mariadb >/dev/null 2>&1
    rm -rf /home/mysql /home/mariadb >/dev/null 2>&1
    
    #Eliminar variables de entorno
    sed -i '/trojan/d' ~/.${SHELL_WAY}rc
    source ~/.${SHELL_WAY}rc

    colorEcho ${GREEN} "Desinstalado éxitosamente!"
}

checkSys() {
    #Comprueba si es Root
    [ $(id -u) != "0" ] && { colorEcho ${RED} "Error: debe ser root para ejecutar este script"; exit 1; }
    if [[ $(uname -m 2> /dev/null) != x86_64 ]]; then
        colorEcho $YELLOW "Ejecute este script en una máquina x86_64."
        exit 1
    fi

    if [[ `command -v apt-get` ]];then
        PACKAGE_MANAGER='apt-get'
    elif [[ `command -v dnf` ]];then
        PACKAGE_MANAGER='dnf'
    elif [[ `command -v yum` ]];then
        PACKAGE_MANAGER='yum'
    else
        colorEcho $RED "No es compatible con el sistema operativo!"
        exit 1
    fi

    # Agregado automáticamente cuando falta la ruta /usr/local/bin
    [[ -z `echo $PATH|grep /usr/local/bin` ]] && { echo 'export PATH=$PATH:/usr/local/bin' >> /etc/bashrc; source /etc/bashrc; }
}

#Instalación de dependencias
installDependent(){
    if [[ ${PACKAGE_MANAGER} == 'dnf' || ${PACKAGE_MANAGER} == 'yum' ]];then
        ${PACKAGE_MANAGER} install socat crontabs bash-completion -y
    else
        ${PACKAGE_MANAGER} update
        ${PACKAGE_MANAGER} install socat cron bash-completion xz-utils -y
    fi
}

setupCron() {
    if [[ `crontab -l 2>/dev/null|grep acme` ]]; then
        if [[ -z `crontab -l 2>/dev/null|grep trojan-web` || `crontab -l 2>/dev/null|grep trojan-web|grep "&"` ]]; then
            #Calcule la hora real del VPS a las 3 a.m., hora de Beijing
            ORIGIN_TIME_ZONE=$(date -R|awk '{printf"%d",$6}')
            LOCAL_TIME_ZONE=${ORIGIN_TIME_ZONE%00}
            BEIJING_ZONE=8
            BEIJING_UPDATE_TIME=3
            DIFF_ZONE=$[$BEIJING_ZONE-$LOCAL_TIME_ZONE]
            LOCAL_TIME=$[$BEIJING_UPDATE_TIME-$DIFF_ZONE]
            if [ $LOCAL_TIME -lt 0 ];then
                LOCAL_TIME=$[24+$LOCAL_TIME]
            elif [ $LOCAL_TIME -ge 24 ];then
                LOCAL_TIME=$[$LOCAL_TIME-24]
            fi
            crontab -l 2>/dev/null|sed '/acme.sh/d' > crontab.txt
            echo "0 ${LOCAL_TIME}"' * * * systemctl stop trojan-web; "/root/.acme.sh"/acme.sh --cron --home "/root/.acme.sh" > /dev/null; systemctl start trojan-web' >> crontab.txt
            crontab crontab.txt
            rm -f crontab.txt
        fi
    fi
}

installTrojan(){
    local SHOW_TIP=0
    if [[ $UPDATE == 1 ]];then
        systemctl stop trojan-web >/dev/null 2>&1
        rm -f /usr/local/bin/trojan
    fi
    LASTEST_VERSION=$(curl -H 'Cache-Control: no-cache' -s "$VERSION_CHECK" | grep 'tag_name' | cut -d\" -f4)
    echo "Descarga del programa de gestión`colorEcho $BLUE $LASTEST_VERSION`versión..."
    curl -L "$DOWNLAOD_URL/$LASTEST_VERSION/trojan" -o /usr/local/bin/trojan
    chmod +x /usr/local/bin/trojan
    if [[ ! -e /etc/systemd/system/trojan-web.service ]];then
        SHOW_TIP=1
        curl -L $SERVICE_URL -o /etc/systemd/system/trojan-web.service
        systemctl daemon-reload
        systemctl enable trojan-web
    fi
    #Commando para completar las variables de entorno
    [[ -z $(grep trojan ~/.${SHELL_WAY}rc) ]] && echo "source <(trojan completion ${SHELL_WAY})" >> ~/.${SHELL_WAY}rc
    source ~/.${SHELL_WAY}rc
    if [[ $UPDATE == 0 ]];then
        colorEcho $GREEN "¡Programa de administración de trojan instalado éxito!\n"
        echo -e "Ejecute comando `colorEcho $BLUE trojan` para administrar trojan\n\n "
        /usr/local/bin/trojan
    else
        if [[ `cat /usr/local/etc/trojan/config.json|grep -w "\"db\""` ]];then
            sed -i "s/\"db\"/\"database\"/g" /usr/local/etc/trojan/config.json
            systemctl restart trojan
        fi
        /usr/local/bin/trojan upgrade db
        if [[ -z `cat /usr/local/etc/trojan/config.json|grep sni` ]];then
            /usr/local/bin/trojan upgrade config
        fi
        systemctl restart trojan-web
        colorEcho $GREEN "Programa de gestión de trojan actualizado con éxito!\n"
    fi
    setupCron
    [[ $SHOW_TIP == 1 ]] && echo "Acceso al navegador'`colorEcho $BLUE https://nombrededominio`'Para gestiónar multiusuarios de trojan en línea"
}

main(){
    [[ ${HELP} == 1 ]] && help && return
    [[ ${REMOVE} == 1 ]] && removeTrojan && return
    [[ $UPDATE == 0 ]] && echo -e "\033[1;32mSe está instalando el programa de administración de trojan.\033[0m" || echo -e "\033[1;33mSe está actualizando el programa de gestión de trojan.\033[0m"
    checkSys
    [[ $UPDATE == 0 ]] && installDependent
    installTrojan
}

main
