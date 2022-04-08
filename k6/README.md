# Mikado K6 Scripts

## Installation
### Load Tests - K6
https://k6.io/docs/getting-started/installation/
#

### Install the Fixture extension - FAKER
https://github.com/szkiba/xk6-faker
# 

## Steps to run
### k6 run [exaple.js]

### Run K6 login load test script
 ```
 k6 run postClasses.js
 ```
### Debug mode
To enable the debug mode put the following parameter into the command line
 ```
--http-debug="full"
 ```

### IMPORTANT NOTE FOR MAC:
If you have problems with the GO environment variables configure them into bash:
`vi $HOME/.bash_profile`

and apply the changes
`source $HOME/.bash_profile`

Run the app from the $GOBIN path if the problem persists with the following command:
`./k6 run /Users/{path}/k6/players.js`
