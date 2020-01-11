#!/bin/sh
if type bash >/dev/null 2>&1; then
    if [ -d /usr/share/bash-completion/completions ]; then
        metr completion bash --out /usr/share/bash-completion/completions/metr
    fi
fi
if type zsh >/dev/null 2>&1; then
    if [ -d /usr/share/zsh/site-functions ]; then
        metr completion zsh --out /usr/share/zsh/site-functions/_metr
    elif [ -d /usr/local/share/zsh/site-functions ]; then
        metr completion zsh --out /usr/local/share/zsh/site-functions/_metr
    fi
fi
