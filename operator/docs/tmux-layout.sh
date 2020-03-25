#!/bin/sh
# reattach to the session 
# tmux a -t kab
tmux new -s kab\; split-window -h \; split-window -v  \;
# new-window 'test' \; 
#C-b "          split vertically (top/bottom)
#C-b %          split horizontally (left/right)

# search in copy-mode
#https://superuser.com/questions/231002/how-can-i-search-within-the-output-buffer-of-a-tmux-shell
#To search in the tmux history buffer for the current window, press Ctrl-b [ to enter copy mode.
#If you're using emacs key bindings (the default), press Ctrl-s then type the string to search for and press Enter. 
#Press n to search for the same string again. 
#Press Shift-n for reverse search. 
#Press Escape twice to exit copy mode. 
#You can use Ctrl-r to search in the reverse direction.