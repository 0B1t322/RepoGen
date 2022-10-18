module github.com/0B1t322/RepoGen/example

go 1.19

require (
	github.com/0B1t322/RepoGen v0.0.2
	github.com/samber/mo v1.5.1
)

require (
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
)

replace (
    github.com/0B1t322/RepoGen => ../.
)