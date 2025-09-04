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

ROOT_PROJECT="$SHELL_FOLDER/../.."

PIPE_NAME="safe-pipe-test-1"

if [ "$(is_darwin_platform)" == "true" ]; then
    TEST_BINARY="$ROOT_PROJECT/release/pipe-for-parallel-darwin-arm64"
elif [ "$(is_windows_platform)" == "true" ]; then
    TEST_BINARY="$ROOT_PROJECT/release/pipe-for-parallel-windows-arm64"
else
    echo ""
    echo "[ERROR]: only support darwin and windows platform"
    echo ""

    exit 1
fi


read_func() {
    COMMAND="\"$TEST_BINARY\" --action=read --pipe=\"$PIPE_NAME\""
    if ! eval "$COMMAND"; then
        echo ""
        echo "[ERROR]: failed to exec pipe-for-parallel read operation"
        echo ""

        exit 1
    fi
}

write_func() {
    local msg=$1
    COMMAND="\"$TEST_BINARY\" --action=write --pipe=\"$PIPE_NAME\" --message=\"$msg\""
    if ! eval "$COMMAND"; then
        echo ""
        echo "[ERROR]: failed to exec pipe-for-parallel write operation"
        echo ""

        exit 1
    fi
}

write_async() {
    CURRENT_INDEX=$1
    for i in {1..100}; do
        write_func "---- yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy: $CURRENT_INDEX -------- $i ------"
    done
}

exit_pipe() {
    COMMAND="\"$TEST_BINARY\" --action=exit --pipe=\"$PIPE_NAME\""
    if ! eval "$COMMAND"; then
        echo ""
        echo "[ERROR]: failed to exec pipe-for-parallel exit operation"
        echo ""

        exit 1
    fi
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

exit_pipe

end_time_1=$(date +%s)
execution_time_1=$((end_time_1 - start_time_1))
echo ""
echo "----------------------------------------"
echo "Execution time [pipe-for-parallel-test]: $execution_time_1 seconds"
echo "----------------------------------------"
echo ""

cd "$OLD_PWD" || exit >/dev/null 2>&1
