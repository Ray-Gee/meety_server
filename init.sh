#!/usr/bin/osascript
tell application "iTerm"
    create window with default profile
    tell current window
        # tabの分割
        tell current session of current tab
            split vertically with default profile
        end tell

        # tab1の処理
        tell current session of current tab
            write text "Univ0903"
            delay 3
            write text "cd /Users/admin/works/meety/server"
            write text "gin -p 3001 -i run main.go"
        end tell

        # tab2の処理
        tell second session of current tab
            write text "Univ0903"
            delay 3
            write text "cd /Users/admin/works/meety/client"
            write text "yarn start"
        end tell
    end tell
end tell