#!/bin/bash

OLD_PWD=$(pwd)
SHELL_FOLDER=$(
    cd "$(dirname "$0")" || exit
    pwd
)
PROJECT_FOLDER=$SHELL_FOLDER/../..
cd "$SHELL_FOLDER" >/dev/null 2>&1 || exit

# shellcheck source=/dev/null
source "$PROJECT_FOLDER/scripts/base/env.sh"

FIFO_PATH="$SHELL_FOLDER/.bash-fifo-test.fifo"
mkfifo "$FIFO_PATH"

read_func() {
    while read -r fifo_data <"$FIFO_PATH"; do
        echo -e "$fifo_data"
    done
}

write_func() {
    local msg=$1

    echo "$msg" >"$FIFO_PATH"
}

write_async() {
    CURRENT_INDEX=$1
    for i in {1..100}; do
        write_func "---- yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy: $CURRENT_INDEX -------- $i ------"
    done
}

close_fifo() {
    rm -rf "$FIFO_PATH"
}

read_func &

start_time_1=$(date +%s)

write_async 1 &
write_pid1=$!
write_async 2 &
write_pid2=$!
write_async 3 &
write_pid3=$!
write_async 4 &
write_pid4=$!
write_async 5 &
write_pid5=$!

wait $write_pid1 $write_pid2 $write_pid3 $write_pid4 $write_pid5

close_fifo

end_time_1=$(date +%s)
execution_time_1=$((end_time_1 - start_time_1))
echo ""
echo "----------------------------------------"
echo "Execution time [bash-fifo-test]: $execution_time_1 seconds"
echo "----------------------------------------"
echo ""

cd "$OLD_PWD" || exit >/dev/null 2>&1
