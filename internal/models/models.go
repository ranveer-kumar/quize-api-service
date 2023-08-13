package models

import "time"

type Quize struct {
	ID           string     `json:"id,omitempty" bson:"_id,omitempty"` //primitive.ObjectID
	Creator      string     `json:"creator,omitempty"  bson:"creator,omitempty"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Schedule     Schedule   `json:"schedule"`
	Questions    []Question `json:"questions"`
	Participants int64      `json:"participants"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	IsPublic     bool       `json:"isPublic"`
	Analytics    Analytics  `json:"analytics"`
	Version      int        `json:"version"`
}

type Analytics struct {
	Views         int64 `json:"views"`
	Participation int64 `json:"participation"`
}

type Question struct {
	Text               string   `json:"text"`
	Choices            []string `json:"choices"`
	CorrectChoiceIndex int64    `json:"correctChoiceIndex"`
}

type Schedule struct {
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	DurationSeconds int64     `json:"durationSeconds"`
}

func (q *Quize) ConvertTimeToIST() {
	istLocation, _ := time.LoadLocation("Asia/Kolkata")
	q.CreatedAt = q.CreatedAt.In(istLocation)
	q.UpdatedAt = q.UpdatedAt.In(istLocation)
	q.Schedule.StartDate = q.Schedule.StartDate.In(istLocation)
	q.Schedule.EndDate = q.Schedule.EndDate.In(istLocation)
}

// FormatCreatedAt returns the CreatedAt field in the specified format
func (q *Quize) FormatCreatedAt() string {
	istLocation, _ := time.LoadLocation("Asia/Kolkata")
	istTime := q.CreatedAt.In(istLocation)
	return istTime.Format(time.RFC3339)
}
