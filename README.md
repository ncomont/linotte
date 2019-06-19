Here you'll find the code of an abandoned project. Not really documented, not really tested neither ; but better here than in a trash.
---------------------


TODO
---------------------

- Write scripts to reset DB, re-ingest, start devenv, start prodenv,
- Create table VernacularGroup and export Verb.FirstLevelVernacularGroup and Verb.SecondLevelVernacularGroup into it,
- Implement console (user creation, etc.),
- Link job result to a Taxon instead of a simple id.



Pull the `backoffice` submodule
---------------------

```
git submodule init
git submodule update
```


Ingest
---------------------

Put ingests files under `docker/volumes/ingest/data`, then start the `db` docker service :

```
./build-docker.sh
cd docker
docker network create linotte-network
docker-compose up db
```


Start the `ingest` docker service :

```
docker-compose -f docker-compose-ingest.yml up
```

When done :

```
docker-compose -f docker-compose-ingest.yml down
```



Development
---------------------

Start the `db` docker service :

```
./build-docker.sh
cd docker
docker network create linotte-network
docker-compose -f docker-compose-dev.yml up db
```

Generate API keys :

```
openssl genrsa -out app.rsa 1024
openssl rsa -in app.rsa -pubout > app.rsa.pub
```

And update settings.

Build and run `Linotte API` on port :

```
go build && ./linotte api -endpoint :10000
```

Start `Linotte backoffice`'s webpack development server :

```
cd backoffice
npm run server
```

Then, go to `http://localhost:8080`.

When done :

```
docker-compose -f docker-compose-ingest.yml down
```



Production
---------------------

Prepare the docker package :

```
./build-docker.sh
```

Copy the `docker` folder where it is needed. See __ingest__ section if necessary. When ready :

```
docker network create linotte-network
docker-compose up
```
