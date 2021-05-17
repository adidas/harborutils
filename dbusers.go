package main

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func FixEmptyEmails(db *gorm.DB) {
	// alternative to be cool
	// 	db.Exec("update harbor_user set email=username where email='' AND  username LIKE='%@%'")

	var users []DbHarborUser
	db.Where("Email = ?", "").Find(&users)
	for _, user := range users {
		fmt.Printf("Fix empty email: %s", user)
		if strings.Contains(user.Username, "@") {
			user.Email = user.Username
			// db.Model(&user).Update("Email", user.Username)
		} else {
			user.Email = fmt.Sprintf("%s@nodamain.io", user.Username)
			// db.Model(&user).Update("Email", username)
		}
		db.Save(&user)
		fmt.Printf("email updated: %s\n", user)

	}
}

// update TargetUser if the data has been modified
func UpdateTargerUser(uSource DbHarborUser, uTarget DbHarborUser, db *gorm.DB) {
	updated := false
	if uSource.UpdateTime != uTarget.UpdateTime {
		// change by a copy using Reflection
		uTarget.Email = uSource.Email
		// uTarget.Password = uSource.Password
		uTarget.Realname = uSource.Realname
		uTarget.Deleted = uSource.Deleted
		uTarget.SysadminFlag = uSource.SysadminFlag
		uTarget.CreationTime = uSource.CreationTime
		uTarget.UpdateTime = uSource.UpdateTime
		db.Save(&uTarget)
		updated = true
	}
	if uSource.OidcUser.Secret != uTarget.OidcUser.Secret {
		uTarget.OidcUser.Secret = uSource.OidcUser.Secret
		db.Save(&uTarget.OidcUser)
		updated = true
	}
	if updated {
		fmt.Printf("User updated: %s\n", uTarget)
	}
}

func CreateNewUser(user DbHarborUser, db *gorm.DB) {
	newUser := DbHarborUser{
		// change by a copy using Reflection
		Username:        user.Username,
		Email:           user.Email,
		Password:        user.Password,
		Realname:        user.Realname,
		Comment:         user.Comment,
		Deleted:         user.Deleted,
		ResetUuid:       user.ResetUuid,
		Salt:            user.Salt,
		SysadminFlag:    user.SysadminFlag,
		CreationTime:    user.CreationTime,
		UpdateTime:      user.UpdateTime,
		PasswordVersion: user.PasswordVersion,
	}

	oidc := DBOidcUser{
		Secret:       user.OidcUser.Secret,
		Subiss:       user.OidcUser.Subiss,
		Token:        user.OidcUser.Token,
		CreationTime: user.OidcUser.CreationTime,
		UpdateTime:   user.OidcUser.UpdateTime,
	}
	newUser.OidcUser = &oidc
	db.Save(&newUser)
	fmt.Printf("User created: %s", newUser)
}

func SyncUsersDatabase(dbSource *gorm.DB, dbTarget *gorm.DB) {
	var usersTarget []DbHarborUser
	var mUsers = make(map[string]DbHarborUser)
	// dbTarget.Preload("OidcUser").Where("Username not in ( 'admin' , 'anonymous')").Find(&usersTarget)
	dbTarget.Joins("OidcUser").Find(&usersTarget, "Username not in ( 'admin' , 'anonymous')")
	for _, user := range usersTarget {
		if user.OidcUser == nil {
			fmt.Printf(" %s no Oidc\n", user)
		} else {
			mUsers[user.Username] = user
		}
	}
	var usersSource []DbHarborUser
	// dbSource.Preload("DbOidcUser").Where("Username not in ( 'admin' , 'anonymous')").Find(&usersSource)
	dbSource.Joins("OidcUser").Find(&usersSource, "Username not in ( 'admin' , 'anonymous')")

	for _, userSource := range usersSource {
		if userTarget, ok := mUsers[userSource.Username]; ok {
			// check if the registry needs to be updated
			UpdateTargerUser(userSource, userTarget, dbTarget)
		} else {
			// create new record
			CreateNewUser(userSource, dbTarget)
		}
	}
}
