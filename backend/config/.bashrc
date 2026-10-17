# ~/.bashrc: executed by bash(1) for non-login shells.

# Note: PS1 and umask are already set in /etc/profile. You should not
# need this unless you want different defaults for root.
# PS1='${debian_chroot:+($debian_chroot)}\h:\w\$ '
# umask 022

# You may uncomment the following lines if you want `ls' to be colorized:
# export LS_OPTIONS='--color=auto'
# eval "`dircolors`"
# alias ls='ls $LS_OPTIONS'
# alias ll='ls $LS_OPTIONS -l'
# alias l='ls $LS_OPTIONS -lA'

alias ll='ls -ltrah --time-style=long-iso --color=auto'

PS1_1='\n\[\033[01;31m\]________________________________________________________'
PS1_2='\n\[\033[01;32m\]`date +%Y-%m-%d\ %H:%M:%S` '
PS1_3='\u@\h \n'
PS1_4='\[\033[01;36m\]`pwd`\[\033[00m\]\n'

export PS1=$PS1_1$PS1_2$PS1_3$PS1_4

