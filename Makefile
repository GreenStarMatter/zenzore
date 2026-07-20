GO_VERSION := 1.18

.PHONY: install-go init-go tf-init tf-plan tf-apply tf-destroy db-start db-stop db-migrate db-setup

setup: install-go init-go
tf-setup: tf-init tf-plan tf-apply
db-setup: db-start db-migrate db-stop

install-go:
	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
	rm go$(GO_VERSION).linux-amd64.tar.gz


init-go:
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}/.bashrc
	echo 'export PATH=$$PATH:$${HOME}:/go/bin' >> $${HOME}/.bashrc

build:
	go build -o zenzore main.go

report:
	go tool cover -html=coverage.out -o cover.html

test:
	go test ./... -coverprofile=coverage.out

coverage:
	go tool cover -func coverage.out | grep "total:" | \
	awk '{print ((int($$3) > 80) != 1) }'

tf-init:
	cd infra && terraform init

tf-plan:
	cd infra && terraform plan

tf-apply:
	cd infra && terraform apply

tf-destroy:
	cd infra && terraform destroy
db-start:
	gcloud sql instances patch zenzore-registry --activation-policy=ALWAYS
db-stop:
	gcloud sql instances patch zenzore-registry --activation-policy=NEVER
db-migrate:
	PGPASSWORD=$(CLOUDSQL_PASSWORD) gcloud sql connect zenzore-registry \
		--user=postgres \
		--database=zenzore_registry < migrations/create_registry_tables.sql
