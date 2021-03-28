module hzsqlcl

go 1.15

require (
	github.com/c-bata/go-prompt v0.2.6
	github.com/candid82/liner v1.4.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/hazelcast/hazelcast-go-client/v4 v4.0.0
)

replace github.com/hazelcast/hazelcast-go-client/v4 v4.0.0 => github.com/yuce/hazelcast-go-client/v4 v4.0.0-temp.2
