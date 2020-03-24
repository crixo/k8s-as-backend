#!/bin/sh
# reattach to the session 
# tmux a -t kab
tmux new -s kab\; split-window -h \; split-window -v  \;
# new-window 'test' \; 
#C-b "          split vertically (top/bottom)
#C-b %          split horizontally (left/right)