#!/bin/sh
if [ -f /usr/share/bash-completion/completions/metr ]; then
    rm /usr/share/bash-completion/completions/metr
fi
if [ -f /usr/share/zsh/site-functions/_metr ]; then
    rm /usr/share/zsh/site-functions/_metr
fi
if [ -f /usr/local/share/zsh/site-functions/_metr ]; then
    rm /usr/local/share/zsh/site-functions/_metr
fi
