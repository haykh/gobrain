#!/usr/bin/env bash

rm -rf /tmp/.gobrain/ &&
  mkdir -p /tmp/.gobrain/ &&
  cp -r template/tasks /tmp/.gobrain/tasks &&
  cp -r template/random_notes /tmp/.gobrain/random_notes &&
  cp -r template/daily_notes /tmp/.gobrain/daily_notes &&
  go build . &&
  ./gobrain -k -d -h /tmp/.gobrain
