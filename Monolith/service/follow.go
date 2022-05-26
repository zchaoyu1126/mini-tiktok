package service

import (
	"fmt"
	"mini-tiktok/repository"
)

func RelationAction(fromUserID, toUserID int64, actionType int64) error {
	if actionType == 1 {
		repository.NewMysqlDaoInstance().Follow(fromUserID, toUserID)

		user1, _ := repository.NewMysqlDaoInstance().QueryUserByID(fromUserID)
		user1.FollowCount++
		repository.NewMysqlDaoInstance().UpdateUser(user1)
		// repository.NewMysqlDaoInstance().ModifyFollowCount(fromUserID, user1.FollowCount+1)

		user2, _ := repository.NewMysqlDaoInstance().QueryUserByID(toUserID)
		repository.NewMysqlDaoInstance().ModifyFollowerCount(toUserID, user2.FollowerCount+1)
		fmt.Println(user1, user2)
	} else if actionType == 2 {
		repository.NewMysqlDaoInstance().DisFollow(fromUserID, toUserID)
	}
	return nil
}
