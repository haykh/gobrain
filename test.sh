#!/usr/bin/env bash

rm -rf ~/.gobrain/tasks/ &&
  cp -r template/tasks ~/.gobrain/tasks &&
  go run .
