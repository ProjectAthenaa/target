buildLocal:
	docker build --build-arg GH_TOKEN=$(token) -e DEBUG=1 -t athena/modules/target:latest .

build:
	docker build --build-arg GH_TOKEN=$(token)  -t registry.digitalocean.com/athenabot/modules/target:latest .

push:
	doctl auth switch --context athena
	doctl registry login
	make build
	docker push registry.digitalocean.com/athenabot/modules/target:latest

rollout:
	doctl kubernetes cluster kubeconfig save athena
	kubectl rollout restart deployments authentication -n general
	kubectl rollout status deployments authentication -n general

deploy:
	echo "Starting deployment"
	make push
	echo "Pushed image to docker"
	make rollout
	echo "Rolled out updates"
	echo ""
	echo "Deployment Done!"