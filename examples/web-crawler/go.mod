module action

go 1.14

replace crawler => ./crawler

require (
	crawler v0.0.0-00010101000000-000000000000
	github.com/apache/openwhisk-client-go v0.0.0-20210313152306-ea317ea2794c
)