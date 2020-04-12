package main

import (
	"fmt"
)

func ExampleIdeDockerRun() {
	*DryRunFlag = true
	DryRunPush("", "172.17.0.2", "641792b3e0112c8fa1896b8944a846dbbab88fe5729f3d464e71475afd9e6057",
		"", "Error:",
		"", "172.17.0.2")
	fmt.Println(1, ideDockerRun("/tmp"))
	fmt.Println(2, ideDockerRun("/tmp"))
	fmt.Println(3, ideDockerRun(""))
	// Output:
	// docker pull iosdk/iosdk-ide:test
	// docker inspect --format={{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}} openwhisk
	// docker run -d -p 3000:3000 --rm --name ide-js --add-host=openwhisk:172.17.0.2 -v /tmp:/home/project iosdk/iosdk-ide:test
	// 1 <nil>
	// docker pull iosdk/iosdk-ide:test
	// docker inspect --format={{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}} openwhisk
	// 2 cannot find openwhisk
	// docker pull iosdk/iosdk-ide:test
	// docker inspect --format={{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}} openwhisk
	// docker run -d -p 3000:3000 --rm --name ide-js --add-host=openwhisk:172.17.0.2 iosdk/iosdk-ide:test
	// 3 <nil>
}

func ExampleIdeDockerRm() {
	*DryRunFlag = true
	DryRunPush("172.17.0.3")
	IdeDestroy()
	// Output:
	// Destroying IDE...
	// docker kill ide-js
	// 172.17.0.3
}
