package main

import (
	"fmt"
	"io/ioutil"
)

func ExampleStart() {
	//*DryRunFlag = true
	fmt.Println("=== Init ===")
	fmt.Println(Start())
	DryRunPush("/tmp/iosdk-test/javascript", "javascript", "123456")
	dir, err := Init("", "", ioutil.Discard)
	fmt.Println(dir, err)
	fmt.Println(Configure(dir))
	fmt.Print(run("ls -a /tmp/iosdk-test/.io*"))
	fmt.Println("=== Start ===")
	DryRunPush(MinDockerVersion, "", "123", "", "1.2.3.4", "", "", "", "172.17.0.2")
	fmt.Println(Start())
	// Output:
	// === Init ===
	// You need to run 'iosdk init ', first.
	// stat /tmp/iosdk-test/.iosdk: no such file or directory
	// Select one of the available templates for importers, or provide your own.
	// The javascript template is for Excel import.
	// The java template is for SQL import.
	// The python template is for REST import.
	// The github template requires a github repo (user/path).
	// Preparing work directory /tmp/iosdk-test/javascript for https://github.com/pagopa/io-sdk-javascript
	// Done.
	// /tmp/iosdk-test/javascript <nil>
	// Wrote /tmp/iosdk-test/.iosdk
	// <nil>
	// /tmp/iosdk-test/.iosdk
	// === Start ===
	// docker version --format {{.Server.Version}}
	// WARNING: using default OpenWhisk key
	// Deploying Redis...
	// docker pull library/redis:5
	// docker run -d -p 6379:6379 --rm --name redis --hostname redis library/redis:5
	//
	// Deploying Whisk...
	// docker pull iosdk/iosdk-openwhisk:test
	// docker inspect --format={{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}} redis
	// docker run -d -p 3280:3280 --rm --name openwhisk --hostname openwhisk -e CONTAINER_EXTRA_ENV=__OW_REDIS=1.2.3.4 -e CONFIG_FORCE_whisk_users_guest=23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP -v //var/run/docker.sock:/var/run/docker.sock iosdk/iosdk-openwhisk:test
	// docker exec openwhisk waitready
	//
	// Deploying IDE...
	// docker pull iosdk/iosdk-ide:test
	// docker inspect --format={{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}} openwhisk
	// docker run -d -p 3000:3000 --rm --name ide-js --add-host=openwhisk:172.17.0.2 -v /tmp/iosdk-test/javascript:/home/project iosdk/iosdk-ide:test
	// <nil>
}

func ExampleStop() {
	*DryRunFlag = true
	DryRunPush()
	fmt.Println(Stop())
	// Output:
	// Destroying IDE...
	// docker kill ide-js
	//
	// Destroying Whisk...
	// docker exec openwhisk stop
	//
	// Destroying Redis...
	// docker stop redis
	//
	// <nil>
}
