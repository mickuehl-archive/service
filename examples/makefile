.PHONY = all site deploy

PROJECT_ID = txsvc-service-example

all: site deploy

site:
	cd web && gatsby build && cd ..
	rm -rf app/public
	cp -R web/public app/public
	
deploy:
	cd app && gcloud -q --project=${PROJECT_ID} app deploy . && cd ..
