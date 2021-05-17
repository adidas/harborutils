package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

// https://registry-dub.tools.3stripes.net/api/v2.0/robots?q=Level%3Dproject%2CProjectID%3D35&page_size=15&page=1

func queryByProyect(project ProjectResponse) string {
	return fmt.Sprint("&q=Level%3Dproject%2CProjectID%3D", project.ProjectId)
}

func _getRobots(harborServer, user, pass, apiVersion, query string) ([]RobotResponse, error) {
	var robots []RobotResponse
	cnt := 1
	for {
		url := fmt.Sprintf("%v/api/%vrobots?page=%v%s",
			harborServer,
			apiVersion,
			cnt,
			query)

		res, body := client(
			ClientPrt{
				Url:         url,
				Method:      "GET",
				ContentType: "application/json",
				User:        user,
				Password:    pass,
			},
		)

		if res.StatusCode < 399 && res.StatusCode > 100 {
			var aux []RobotResponse
			json.Unmarshal([]byte(body), &aux)
			robots = append(robots, aux...)
		} else {
			log.Fatal("Error getting robots")
		}

		l := res.Header.Get("Link")
		if strings.Contains(l, "next") {
			cnt++
		} else {
			break
		}
	}
	if len(robots) > 0 {
		return robots, nil
	}
	return robots, errors.New("Error getting robots")
}

func getRobots(harborServer, user, pass, apiVersion string) []RobotResponse {
	// getPublicRobots
	robots, _ := _getRobots(harborServer, user, pass, apiVersion, "")

	// getProyects
	projects := listProjects(harborServer, user, pass, apiVersion)
	// var proyectsRobots []RobotResponse
	for _, project := range projects {
		aux, _ := _getRobots(harborServer, user, pass, apiVersion, queryByProyect(project))
		robots = append(robots, aux...)
	}
	return robots
}

func updateRobotTarget(harborTarget, userTarget, passTarget, apiVersionTarget string, robotSource, robotTarget RobotResponse) {
	url := fmt.Sprintf("%v/api/%vrobots/%d",
		harborTarget,
		apiVersionTarget,
		robotTarget.ID,
	)
	id := robotSource.ID
	robotSource.ID = robotTarget.ID

	res, _ := client(
		ClientPrt{
			Url:         url,
			Method:      "PUT",
			ContentType: "application/json",
			User:        userTarget,
			Password:    passTarget,
			Body:        robotSource,
		},
	)
	robotSource.ID = id
	if res.StatusCode < 100 || (res.StatusCode > 399 && res.StatusCode != 409) {
		// print status code because some groups might not be valid anymore and import should continue
		fmt.Printf("Status Code: %d while trying to update robot, robotSource: %s, robotTarget: %s\n", res.StatusCode, robotSource, robotTarget)
	} else {
		fmt.Printf("Robot updated: %s\n", robotTarget)
	}
}

func updateRobotDb(dbSource, dbTarget *gorm.DB, robotSourceId, robotTargetId int) {
	var dbRobotSource, dbRobotTarget DbRobot
	dbSource.First(&dbRobotSource, robotSourceId)
	dbTarget.First(&dbRobotTarget, robotTargetId)
	if dbRobotTarget.Secret != dbRobotSource.Secret || dbRobotTarget.Salt != dbRobotSource.Salt {
		dbRobotTarget.Secret = dbRobotSource.Secret
		dbRobotTarget.Salt = dbRobotSource.Salt
		dbTarget.Save(&dbRobotTarget)
		fmt.Printf("UpdateRobot secret, robot_source: %s robot_target: %s\n", dbRobotSource, dbRobotTarget)
	}

}

func createRobotTarget(harborTarget, userTarget, passTarget, apiVersionTarget string, robot RobotResponse) RobotCreatedReponse {
	url := fmt.Sprintf("%v/api/%vrobots",
		harborTarget,
		apiVersionTarget,
	)

	name := robot.Name
	if robot.Level == "system" {
		robot.Name = strings.Split(robot.Name, "$")[1]
	} else if robot.Level == "project" {
		robot.Name = strings.Split(robot.Name, "+")[1]
	} else {
		fmt.Printf("Error when try identified robot type. Robot: %s\n", robot)
	}

	id := robot.ID
	robot.ID = 0

	res, body := client(
		ClientPrt{
			Url:         url,
			Method:      "POST",
			ContentType: "application/json",
			User:        userTarget,
			Password:    passTarget,
			Body:        robot,
		},
	)
	robot.ID = id
	robot.Name = name
	var robotResponse RobotCreatedReponse
	if res.StatusCode < 100 || (res.StatusCode > 399 && res.StatusCode != 409) {
		// print status code because some groups might not be valid anymore and import should continue
		fmt.Printf("Status Code: %d while trying to create robot: %s\n", res.StatusCode, robot)
	} else {
		fmt.Printf("Robot created: %s\n", robot)
		json.Unmarshal([]byte(body), &robotResponse)
	}
	return robotResponse
}

func deleteRobot(harborServer, user, pass, apiVersion string, robot RobotResponse) {
	url := fmt.Sprintf("%v/api/%vrobots/%d",
		harborServer,
		apiVersion,
		robot.ID)

	res, _ := client(
		ClientPrt{
			Url:         url,
			Method:      "DELETE",
			ContentType: "application/json",
			User:        user,
			Password:    pass,
		},
	)
	if res.StatusCode < 100 || (res.StatusCode > 399 && res.StatusCode != 409) {
		// print status code because some groups might not be valid anymore and import should continue
		fmt.Printf("Status Code: %d while trying to delete robot: %s", res.StatusCode, robot)
	} else {
		fmt.Printf("Robot %s deleted: ", robot)
	}

}

func syncRobots(harborSource, userSource, passSource, apiVersionSource, harborTarget, userTarget, passTarget, apiVersionTarget string, dbSource, dbTarget *gorm.DB) {
	robotsSource := getRobots(harborSource, userSource, passSource, apiVersionSource)
	robotsTarget := getRobots(harborTarget, userTarget, passTarget, apiVersionSource)
	mRobotsTarget := make(map[string]RobotResponse)
	mRobotsSource := make(map[string]RobotResponse)

	for _, robot := range robotsTarget {
		mRobotsTarget[robot.Name] = robot
	}

	for _, robot := range robotsSource {
		mRobotsSource[robot.Name] = robot
	}

	for _, robotSource := range robotsSource {
		// robot account is not legacy
		if robotSource.Editable {
			if robotTarget, ok := mRobotsTarget[robotSource.Name]; ok {
				updateRobotTarget(harborTarget, userTarget, passTarget, apiVersionTarget, robotSource, robotTarget)
				updateRobotDb(dbSource, dbTarget, robotSource.ID, robotTarget.ID)

			} else {
				robotCreated := createRobotTarget(harborTarget, userTarget, passTarget, apiVersionTarget, robotSource)
				if robotCreated.ID != 0 {
					updateRobotDb(dbSource, dbTarget, robotSource.ID, robotCreated.ID)
				}
			}
		}
	}

	for _, robot := range robotsTarget {
		if _, ok := mRobotsSource[robot.Name]; !ok {
			deleteRobot(harborTarget, userTarget, passTarget, apiVersionTarget, robot)
		}
	}

}
