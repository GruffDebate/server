package gruff

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	setupDB()
	defer teardownDB()

	user := User{
		Name:     fmt.Sprintf("Imma User"),
		Username: "ImmaUser",
		Email:    "immauser@gruff.org",
		Password: "monkey",
		Image:    "https://i.ytimg.com/vi/hYuViV9NgzA/hqdefault.jpg",
		Curator:  true,
		Admin:    true,
		URL:      "https://thetruth2020.org/",
	}

	saved := User{}
	saved.Username = user.Username
	err := saved.Load(CTX)
	assert.Empty(t, saved.Key)

	err = user.Create(CTX)
	assert.NoError(t, err)
	saved = User{}
	saved.Username = user.Username
	err = saved.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Nil(t, saved.DeletedAt)
	assert.Equal(t, "", saved.Password)
	assert.NotEmpty(t, saved.HashedPassword)
	assert.Equal(t, user.Name, saved.Name)
	assert.Equal(t, user.Username, saved.Username)
	assert.Equal(t, user.Email, saved.Email)
	assert.Equal(t, user.Image, saved.Image)
	assert.True(t, saved.Curator)
	assert.True(t, saved.Admin)
	assert.Equal(t, user.URL, saved.URL)
}

func TestUserScoreFor(t *testing.T) {
	setupDB()
	defer teardownDB()

	u := User{
		Username: "TheBigScore",
	}
	err := u.Create(CTX)
	assert.NoError(t, err)

	claim := Claim{
		Title: "Dude, I totally scored!",
	}
	err = claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Scored? Like, left scratch marks?",
	}
	err = arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	startTime := time.Now()

	err = u.Score(CTX, &claim, 0.90)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	cScore, err := claim.Score(CTX)
	assert.NoError(t, err)
	aScore, err := arg.Score(CTX)
	assert.NoError(t, err)
	aStr, err := arg.Strength(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.90), claim.Truth)
	assert.Equal(t, float32(1.00), arg.Relevance)
	assert.Equal(t, float32(0.50), arg.Str)
	assert.Equal(t, float32(0.90), cScore)
	assert.Equal(t, float32(1.00), aScore)
	assert.Equal(t, float32(0.50), aStr)

	score, err := u.ScoreFor(CTX, &claim)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.90), score.Score)

	score, err = u.ScoreFor(CTX, &arg)
	assert.NoError(t, err)
	assert.Nil(t, score)

	claim.QueryAt = &startTime
	score, err = u.ScoreFor(CTX, &claim)
	assert.NoError(t, err)
	assert.Nil(t, score)
	cScore, err = claim.Score(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.50), cScore)
	claim.QueryAt = nil

	err = u.Score(CTX, &arg, 0.75)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	firstScoresTime := time.Now()

	cScore, err = claim.Score(CTX)
	assert.NoError(t, err)
	aScore, err = arg.Score(CTX)
	assert.NoError(t, err)
	aStr, err = arg.Strength(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.90), claim.Truth)
	assert.Equal(t, float32(0.75), arg.Relevance)
	assert.Equal(t, float32(0.375), arg.Str)
	assert.Equal(t, float32(0.90), cScore)
	assert.Equal(t, float32(0.75), aScore)
	assert.Equal(t, float32(0.375), aStr)

	score, err = u.ScoreFor(CTX, &claim)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.90), score.Score)

	score, err = u.ScoreFor(CTX, &arg)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.75), score.Score)

	arg.QueryAt = &startTime
	score, err = u.ScoreFor(CTX, &arg)
	assert.NoError(t, err)
	assert.Nil(t, score)
	aScore, err = arg.Score(CTX)
	assert.NoError(t, err)
	aStr, err = arg.Strength(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(1.00), aScore)
	assert.Equal(t, float32(0.50), aStr)
	arg.QueryAt = nil

	err = u.Score(CTX, &claim, 0.10)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = u.Score(CTX, &arg, 0.35)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	secondScoresTime := time.Now()

	cScore, err = claim.Score(CTX)
	assert.NoError(t, err)
	aScore, err = arg.Score(CTX)
	assert.NoError(t, err)
	aStr, err = arg.Strength(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.10), claim.Truth)
	assert.Equal(t, float32(0.35), arg.Relevance)
	assert.Equal(t, float32(0.175), arg.Str)
	assert.Equal(t, float32(0.10), cScore)
	assert.Equal(t, float32(0.35), aScore)
	assert.Equal(t, float32(0.175), aStr)

	score, err = u.ScoreFor(CTX, &claim)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.10), score.Score)

	score, err = u.ScoreFor(CTX, &arg)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.35), score.Score)

	claim.QueryAt = &firstScoresTime
	score, err = u.ScoreFor(CTX, &claim)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.90), score.Score)
	cScore, err = claim.Score(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.90), cScore)
	claim.QueryAt = nil

	arg.QueryAt = &firstScoresTime
	score, err = u.ScoreFor(CTX, &arg)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.75), score.Score)
	aScore, err = arg.Score(CTX)
	assert.NoError(t, err)
	aStr, err = arg.Strength(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.75), aScore)
	assert.Equal(t, float32(0.375), aStr)
	arg.QueryAt = nil

	err = claim.Update(CTX, Updates{})
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = arg.Update(CTX, Updates{})
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim.Load(CTX)
	assert.NoError(t, err)
	err = arg.Load(CTX)
	assert.NoError(t, err)

	cScore, err = claim.Score(CTX)
	assert.NoError(t, err)
	aScore, err = arg.Score(CTX)
	assert.NoError(t, err)
	aStr, err = arg.Strength(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.10), claim.Truth)
	assert.Equal(t, float32(0.35), arg.Relevance)
	assert.Equal(t, float32(0.175), arg.Str)
	assert.Equal(t, float32(0.10), cScore)
	assert.Equal(t, float32(0.35), aScore)
	assert.Equal(t, float32(0.175), aStr)

	score, err = u.ScoreFor(CTX, &claim)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.10), score.Score)

	score, err = u.ScoreFor(CTX, &arg)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.35), score.Score)

	claim.QueryAt = &secondScoresTime
	score, err = u.ScoreFor(CTX, &claim)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.10), score.Score)
	cScore, err = claim.Score(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.10), cScore)
	claim.QueryAt = nil

	arg.QueryAt = &secondScoresTime
	score, err = u.ScoreFor(CTX, &arg)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.35), score.Score)
	aScore, err = arg.Score(CTX)
	assert.NoError(t, err)
	aStr, err = arg.Strength(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.35), aScore)
	assert.Equal(t, float32(0.175), aStr)
	arg.QueryAt = nil

	err = u.Score(CTX, &claim, 0.50)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = u.Score(CTX, &arg, 0.55)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	cScore, err = claim.Score(CTX)
	assert.NoError(t, err)
	aScore, err = arg.Score(CTX)
	assert.NoError(t, err)
	aStr, err = arg.Strength(CTX)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.50), claim.Truth)
	assert.Equal(t, float32(0.55), arg.Relevance)
	assert.Equal(t, float32(0.275), arg.Str)
	assert.Equal(t, float32(0.50), cScore)
	assert.Equal(t, float32(0.55), aScore)
	assert.Equal(t, float32(0.275), aStr)

	score, err = u.ScoreFor(CTX, &claim)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.50), score.Score)

	score, err = u.ScoreFor(CTX, &arg)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.55), score.Score)
}

// TODO: test update
// TODO: test change password
// TODO: test validations
