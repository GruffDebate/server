# gruff-server
Server-side code for the open source "Wikipedia for Debates", called Gruff

This project is the server-side API written in Golang, with ArangoDB as the backend.

This project is one of many experiments in trying to achieve the goals laid out in The Canonical Debate white paper: https://github.com/canonical-debate-lab/paper

## Database
## Install ArangoDB
If you haven't already, you need to install the [ArangoDB](https://arangodb.com/) database, version 3.1 or later.

## Install Golang
This project uses the new default dependency management tool that has been available since Golang version 1.11. If you haven't already, go to the [Go Programming Language install page](https://golang.org/doc/install) for instructions on how to get it set up in your environment.

## Install ArangoMiGO
ArangoDB schema creation and migration is managed via the [ArangoMiGO](https://github.com/deusdat/arangomigo) tool.

## Set up this project
If you haven't already, clone this project to a local folder.

## Create the database
First, you must create a local configuration file, as required by ArangoMiGO. An example is located in the file `migrations/config.example` of this project. You can make a copy, and then edit the copy to set your own variables.

```bash
cp migrations/config.example migrations/config
echo "I will use the best editor"
emacs migrations/config
```

At a minimum, you will need to customize the `migrationspath` configuration, and the `username` and `password` in `extras`. You should also change the ArangoDB `username` and `password` for the root user if you didn't use the default configurations.

Next, make sure your ArangoDB server is running:

```bash
/usr/local/opt/arangodb/sbin/arangod
```

Finally, run the ArangoMiGO utility to install the database collections and edges.

```bash
arangomigo migrations/config
```

## Optional: Import Data
The database and scheme are the same used by the Canonical Debate Lab's arango-importer project: https://github.com/canonical-debate-lab/arango-importer

If you wish to use the data set provided by that project, please download that project and follow the instructions on how to set up the database and import the sample data. However, it is important to note that the arango-importer project creates a database called "canonical-debate" and this project uses a database called "gruff". TODO: You must change the arango-importer configuration file to point to your "gruff" database in order to import the data for this project.

Warning; the two schemas may get out of sync. In order to guarantee compatibility, specific versions of this project and the arango-importer project will be tagged to show compatibility.

## Developers
If you wish to help with development for this project, then you should set up the test database as well. After following the steps above, you should use the arango shell to create the test database:

```bash
arangosh
127.0.0.1:8529@_system> db._create("gruff_test");
127.0.0.1:8529@_system> exit
```

Next, you should run ArangoMigo again in order to create the schema for the test database, using a separate configuration file:

```bash
cp migrations/config.test.example migrations/config.test
echo "Emacs really is the best editor"
emacs migrations/config.test
...
arangomigo migrations/config.test
```

This should set up your test database, which will be used when running the test suite for this project. Each time you update your code base to the latest version, be sure to run the migration again to make sure your test DB will still be compatible.
