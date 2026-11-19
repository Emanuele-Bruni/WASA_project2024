#!/usr/bin/env sh
set -eu

project_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
suite=${1:-all}

run_backend() {
	if command -v go >/dev/null 2>&1; then
		(cd "$project_dir" && go test ./...)
	elif command -v docker >/dev/null 2>&1; then
		docker run --rm -v "$project_dir:/src" -w /src golang:1.27.3 go test ./...
	else
		echo "error: Go or Docker is required to run backend tests" >&2
		exit 1
	fi
}

run_frontend() {
	(
		cd "$project_dir/webui"
		env YARN_GLOBAL_FOLDER=/tmp/wasa-yarn-global \
			YARN_ENABLE_GLOBAL_CACHE=false \
			node .yarn/releases/yarn-4.5.0.cjs test
	)
}

case "$suite" in
	all)
		run_backend
		run_frontend
		;;
	backend)
		run_backend
		;;
	frontend)
		run_frontend
		;;
	*)
		echo "usage: $0 [all|backend|frontend]" >&2
		exit 2
		;;
esac
