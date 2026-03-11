#!/bin/sh
set -eu

COMMAND="${1:-}"
APP_NAME="${APP_NAME:-platform-truth-mcp}"
BUILD_BIN="${BUILD_BIN:?BUILD_BIN is required}"
CONFIG="${CONFIG:?CONFIG is required}"
RUNTIME_DIR="${RUNTIME_DIR:?RUNTIME_DIR is required}"
PID_FILE="${PID_FILE:?PID_FILE is required}"
CHILD_PID_FILE="${CHILD_PID_FILE:?CHILD_PID_FILE is required}"
REQUEST_PIPE="${REQUEST_PIPE:?REQUEST_PIPE is required}"
LOG_FILE="${LOG_FILE:?LOG_FILE is required}"

is_pid_running() {
  pid_file="$1"
  if [ ! -f "$pid_file" ]; then
    return 1
  fi
  pid="$(cat "$pid_file" 2>/dev/null || true)"
  if [ -z "$pid" ]; then
    return 1
  fi
  kill -0 "$pid" 2>/dev/null
}

cleanup_runtime() {
  rm -f "$CHILD_PID_FILE" "$REQUEST_PIPE"
}

run_supervisor() {
  trap 'if is_pid_running "$CHILD_PID_FILE"; then kill "$(cat "$CHILD_PID_FILE")" 2>/dev/null || true; wait "$(cat "$CHILD_PID_FILE")" 2>/dev/null || true; fi; cleanup_runtime; rm -f "$PID_FILE"' INT TERM EXIT

  exec 3<>"$REQUEST_PIPE"
  "$BUILD_BIN" serve --config "$CONFIG" <&3 >>"$LOG_FILE" 2>&1 &
  child_pid=$!
  echo "$child_pid" > "$CHILD_PID_FILE"
  wait "$child_pid"
}

start_cmd() {
  mkdir -p "$RUNTIME_DIR"

  if [ ! -x "$BUILD_BIN" ]; then
    echo "Binary not found or not executable: $BUILD_BIN" >&2
    exit 1
  fi

  if [ ! -f "$CONFIG" ]; then
    echo "Config not found: $CONFIG" >&2
    exit 1
  fi

  if is_pid_running "$PID_FILE"; then
    echo "$APP_NAME is already running with PID $(cat "$PID_FILE")"
    echo "Request pipe: $REQUEST_PIPE"
    echo "Log file: $LOG_FILE"
    exit 0
  fi

  cleanup_runtime
  rm -f "$PID_FILE"
  mkfifo "$REQUEST_PIPE"
  : > "$LOG_FILE"

  "$0" _run &
  supervisor_pid=$!
  echo "$supervisor_pid" > "$PID_FILE"

  sleep 1
  if is_pid_running "$PID_FILE"; then
    echo "Started $APP_NAME"
    echo "Supervisor PID: $supervisor_pid"
    echo "Request pipe: $REQUEST_PIPE"
    echo "Log file: $LOG_FILE"
    exit 0
  fi

  echo "Failed to start $APP_NAME. Check $LOG_FILE" >&2
  cleanup_runtime
  rm -f "$PID_FILE"
  exit 1
}

status_cmd() {
  if is_pid_running "$PID_FILE"; then
    echo "$APP_NAME is running"
    echo "Supervisor PID: $(cat "$PID_FILE")"
    if is_pid_running "$CHILD_PID_FILE"; then
      echo "Server PID: $(cat "$CHILD_PID_FILE")"
    else
      echo "Server PID: unavailable"
    fi
    echo "Config: $CONFIG"
    echo "Request pipe: $REQUEST_PIPE"
    echo "Log file: $LOG_FILE"
    exit 0
  fi

  echo "$APP_NAME is not running"
  if [ -f "$PID_FILE" ]; then
    echo "Found stale PID file at $PID_FILE"
  fi
  exit 1
}

stop_cmd() {
  if ! is_pid_running "$PID_FILE"; then
    echo "$APP_NAME is not running"
    cleanup_runtime
    rm -f "$PID_FILE"
    exit 0
  fi

  supervisor_pid="$(cat "$PID_FILE")"
  kill "$supervisor_pid" 2>/dev/null || true
  sleep 1

  if kill -0 "$supervisor_pid" 2>/dev/null; then
    kill -9 "$supervisor_pid" 2>/dev/null || true
  fi

  cleanup_runtime
  rm -f "$PID_FILE"
  echo "Stopped $APP_NAME"
}

case "$COMMAND" in
  start)
    start_cmd
    ;;
  status)
    status_cmd
    ;;
  stop)
    stop_cmd
    ;;
  _run)
    run_supervisor
    ;;
  *)
    echo "Usage: $0 {start|status|stop}" >&2
    exit 1
    ;;
esac
