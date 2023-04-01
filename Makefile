.PHONY: deploy
deploy: 
	make -j1 build
	make -j3 deploy-only

.PHONY: build
build: 
	spin build

.PHONY: deploy-only
deploy-only: main preview1 preview2

.PHONY: main
main:
	spin deploy -f spin.toml

.PHONY: preview1
preview1:
	spin deploy -f spin-preview1.toml

.PHONY: preview2
preview2:
	spin deploy -f spin-preview2.toml
