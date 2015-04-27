package user

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func bb_addUser(name, gecos, homedir, passwd string, noCreateHome bool, groups []string, system bool) error {
	err := bb_removeIfUserGroupExists(name)
	if err != nil {
		return err
	}
	userAddArgs := []string{}
	if gecos != "" {
		userAddArgs = append(userAddArgs, "-g", gecos)
	}
	userAddArgs = append(userAddArgs, "-D")
	if !noCreateHome && homedir != "" {
		userAddArgs = append(userAddArgs, "-h", homedir)
	} else {
		userAddArgs = append(userAddArgs, "-H")
	}
	if system {
		userAddArgs = append(userAddArgs, "-S")
	}
	for _, g := range groups {
		userAddArgs = append(userAddArgs, "-G", g)
	}
	userAddArgs = append(userAddArgs, name)
	log.Infof("%+v", userAddArgs)
	cmd := exec.Command("adduser", userAddArgs...)
	if err := cmd.Run(); err != nil {
		return err
	}
	if passwd == "" {
		return nil
	}
	return bb_chPasswd(name, passwd)
}

func bb_chPasswd(user, passwd string) error {
	cmd := exec.Command("chpasswd", "-e")
	cmd.Stdin = strings.NewReader(user + ":" + passwd)
	return cmd.Run()
}

func bb_delUser(name string) error {
	cmd := exec.Command("deluser", name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func bb_delGroup(name string) error {
	cmd := exec.Command("delgroup", name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func bb_removeIfUserGroupExists(name string) error {
	passwdFile, err := os.Open("/etc/passwd")
	if err != nil {
		return err
	}
	passwdScanner := bufio.NewScanner(passwdFile)
	for passwdScanner.Scan() {
		if strings.Contains(passwdScanner.Text(), name+":") {
			log.Infof("user with name=%s found, deleting and re-adding with new config", name)
			if err := bb_delUser(name); err != nil {
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
			if err := bb_delGroup(name); err != nil {
				return err
			}
		}
	}
	return nil
}
