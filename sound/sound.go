package sound

import (
	"fmt"
	"os/exec"
)

/*
	/usr/share/sounds/gnome/default/alerts
	play drip.ogg

*/

func PlayDrip(count int) {
	args := []string{
		"/usr/share/sounds/gnome/default/alerts/drip.ogg",
		"repeat",
		fmt.Sprintf("%d", count),
	}

	exec.Command("play", args...).Start()

}
