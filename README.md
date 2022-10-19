# Currency App

There are some requisites to run the application
- [Go](https://golang.org/doc/install) 1.19
- [Docker](https://docs.docker.com/engine/install/)
- [Make](https://www.gnu.org/software/make/)

# Lint the Challenge
* The command to lint a specific package is the following:
  ```bash
    $ make lint pkg=routes
  ```

* Hovewer,  if you want to lint all packages, then run:
  ```bash
    $ make lint
  ```

# Test the Challenge

There are some commands to run, and test the challenges:

* If you want to test a particular test case then use the following example:
  ```bash
    $ make test pkg=./routes name=TestCompanyRoute_WithWrongLegacyHeaders | tee log.json
  ```
  Note~> you can see the `test` command is the main of this and then the `pkg` tells what package it would be tested, and finally the name of the test case. If you want to save the test result you can add the pipe and the file name as the above example.

* If you want to run all test cases just run the following command:
  ```bash
    $ make test | tee log.json
  ```

# Build and Run
Before to build the project is necessary have docker runing, there are some ways to run the app, the first one is to run the docker-compose.yaml file, which consist in the following steps:

* First make the docker image of the application with the following command:
  ```
  $ make docker-build service=currency version=1.0.0-test
  ```
  Note: you can use the version you want.

* Second is re-write the docker-compose in the app section to put the version you have created in the previous point, then run the docker-compose:
  ```
  $ docker-compose up 
  ```
  Note: you can run it using the -d argument to run it in background mode.

  You might see the app is starting before the database, no worry about it, this is because we are running some sql queries, but be sure start using the app until the database process finish.

  Well, that's it...
