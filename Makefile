.PHONY: test test-backend test-frontend test-browser lint format check

test:
	./test.sh

test-backend:
	./test.sh backend

test-frontend:
	./test.sh frontend

test-browser:
	docker compose -p wasa-e2e -f compose.test.yml build
	docker compose -p wasa-e2e -f compose.test.yml run --rm browser; result=$$?; docker compose -p wasa-e2e -f compose.test.yml down --volumes; exit $$result

lint:
	cd webui && node .yarn/releases/yarn-4.5.0.cjs lint
	cd webui && node .yarn/releases/yarn-4.5.0.cjs format:check

format:
	gofmt -w cmd service
	cd webui && node .yarn/releases/yarn-4.5.0.cjs format

check: test lint
	cd webui && node .yarn/releases/yarn-4.5.0.cjs build-prod
