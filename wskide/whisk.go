package wskide

import (
	"fmt"
)

// WhiskDeploy deploys openwhisk standalone
func WhiskDeploy() error {
	fmt.Println("Deploying Whisk...")
	fmt.Println(whiskDockerRun())
	return nil
}

// WhiskDestroy destroys openwhisk standalone
func WhiskDestroy() error {
	fmt.Println("Destroying Whisk...")
	fmt.Println(Sys("docker exec openwhisk stop"))
	return nil
}

// return empty string if ok, otherwise the error
func whiskDockerRun() string {
	Config, _ := LoadConfig()
	err := Run("docker pull " + OpenwhiskStandaloneImage)
	if err != nil {
		return "cannot pull " + OpenwhiskStandaloneImage
	}
	redisIP := dockerIP("redis")
	if redisIP == nil {
		return "cannot locate redis"
	}
	cmd := fmt.Sprintf(`docker run -d -p 3280:3280
--rm --name openwhisk --hostname openwhisk
-e CONTAINER_EXTRA_ENV=__OW_REDIS=%s -e CONFIG_FORCE_whisk_users_guest=%s
-v //var/run/docker.sock:/var/run/docker.sock %s`, *redisIP, Config.WhiskAPIKey, OpenwhiskStandaloneImage)
	_, err = SysErr(cmd)
	if err != nil {
		return "cannot start server: " + err.Error()
	}

	err = Run("docker exec openwhisk wsk property set apihost http://localhost:3233 --apihost http://localhost:3233 auth " + Config.WhiskAPIKey + " --auth " + Config.WhiskAPIKey)
	if err != nil {
		return "cannot update properties: " + err.Error()
	}

	err = Run("docker exec openwhisk waitready")
	if err != nil {
		return "server readyness error: " + err.Error()
	}

	err = Run("docker cp " + configFile + " openwhisk:/tmp/.iosdk")
	if err != nil {
		return "Error coping config file: " + err.Error()
	}

	err = Run("docker exec openwhisk wsk package update iosdk -P /tmp/.iosdk")
	if err != nil {
		return "error updating iosdk package: " + err.Error()
	}

	return ""

}
