# getOrgName Cli app
cli app to check owner organization name for any web site
to run app from terminal, you should have go install in your OS
to get help


`go run main.go --help`


example


`go run main.go -domains=google.com,github.com`


output
`
    [*] google.com 172.217.20.110 Google
    [*] github.com 192.30.253.112 GitHub
    [*] github.com 192.30.253.113 GitHub
`

to build app

`go build -o output .`