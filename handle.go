package function

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"time"
)

type Questions struct {
	SessionId string
	Question1 string
}

type Score struct {
	SessionId  string
	Time       time.Time
	Level      string
	LevelScore int
}

type GameTime struct{
	GameTimeId string
	SessionId string
	Level string
	Type string
	Time      time.Time
}

var redisHost = os.Getenv("REDIS_HOST")// This should include the port which is most of the time 6379
var redisPassword = os.Getenv("REDIS_PASSWORD")

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	points := 10
	var q Questions

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(req.Body).Decode(&q)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	//Evaluate question answer
	if q.Question1 != "" {
		points = 10
	} else {
		points = -1
	}

	var score Score
	score.Level = "level-2"
	score.LevelScore = points
	score.SessionId = q.SessionId
	score.Time = time.Now()
	scoreJson, err := json.Marshal(score)
	if err != nil {
		fmt.Println(err)
	}
	err = client.RPush("score-"+q.SessionId, string(scoreJson)).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		fmt.Println(err)
	}

	gt := GameTime{
		GameTimeId: "time-" + score.SessionId,
		SessionId:  score.SessionId,
		Level:      score.Level,
		Type:       "end",
		Time:       score.Time,
	}

	gameTimeJson, err := json.Marshal(gt)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = client.RPush(gt.GameTimeId, string(gameTimeJson)).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}


	fmt.Fprintln(res, string(scoreJson))

}
