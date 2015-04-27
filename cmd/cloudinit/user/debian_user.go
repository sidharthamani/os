package user

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func deb_addUser(name, gecos, homedir, passwd string, noCreateHome bool, groups []string, system bool) error {
	err := deb_removeIfUserGroupExists(name)
	if err != nil {
		return err
	}
	userAddArgs := []string{}
	if gecos != "" {
		userAddArgs = append(userAddArgs, "--gecos", gecos)
	}
	userAddArgs = append(userAddArgs, "--disabled-password")
	if !noCreateHome && homedir != "" {
		userAddArgs = append(userAddArgs, "--home", homedir)
	} else {
		userAddArgs = append(userAddArgs, "--no-create-home")
	}
	if system {
		userAddArgs = append(userAddArgs, "--system")
	}
	for _, g := range groups {
		userAddArgs = append(userAddArgs, "--ingroup", g)
	}
	userAddArgs = append(userAddArgs, name)
	cmd := exec.Command("adduser", userAddArgs...)
	if err := cmd.Run(); err != nil {
		return err
	}
	if passwd == "" {
		return nil
	}
	return deb_chPasswd(name, passwd)
}

func deb_chPasswd(user, passwd string) error {
	cmd := exec.Command("chpasswd", "-e")
	cmd.Stdin = strings.NewReader(user + ":" + passwd)
	return cmd.Run()
}

func deb_delUser(name string) error {
	cmd := exec.Command("deluser", name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func deb_delGroup(name string) error {
	cmd := exec.Command("delgroup", name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func deb_removeIfUserGroupExists(name string) error {
	passwdFile, err := os.Open("/etc/passwd")
	if err != nil {
		return err
	}
	passwdScanner := bufio.NewScanner(passwdFile)
	for passwdScanner.Scan() {
		if strings.Contains(passwdScanner.Text(), name+":") {
			log.Infof("user with name=%s found, deleting and re-adding with new config", name)
			if err := deb_delUser(name); err != nil {
				return err
			}
		}
	}
	group, err := os.Open("/etc/group")
	if err != nil {
		return err
	}
	groupsScanner := bufio.NewScanner(group)
	for groupsScanner.Scan() {
		if strings.Contains(groupsScanner.Text(), name+":") {
			log.Infof("group with name=%s found, deleting and re-adding with new config", name)
			if err := deb_delGroup(name); err != nil {
				return err
			}
		}
	}
	return nil
}
