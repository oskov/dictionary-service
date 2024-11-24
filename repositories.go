package main

import "time"

type Word struct {
	ID   int64  `db:"id"`
	Word string `db:"word"`
}

type WordDefinition struct {
	ID         int64  `db:"id"`
	WordID     int64  `db:"wordId"`
	Definition string `db:"definition"`
	Priority   int64  `db:"priority"`
}

type WordDefinitionToUser struct {
	ID               int64 `db:"id"`
	WordDefinitionID int64 `db:"wordDefinitionId"`
	UserID           int64 `db:"userId"`
}

type User struct {
	ID int64 `db:"id"`
}

type SpacedRepetetionCard struct {
	ID                     int64     `db:"id"`
	WordDefenitionToUserID int64     `db:"word_definition_to_user_id"`
	EaseFactor             float64   `db:"ease_factor"` // default 2.5, 5 - easy, 1 - hard
	Interval               int64     `db:"interval"`
	Repetition             int64     `db:"repetition"`
	DueDate                time.Time `db:"due_date"`
}

// Function to calculate the next review date and adjust the interval/ease factor
func reviewCard(card *SpacedRepetetionCard, userAnswerQuality int) {
	// Adjust the ease factor based on the user answer (1-5)
	if userAnswerQuality < 3 {
		// Reset repetition count on failure
		card.Repetition = 0
		card.Interval = 1 // Show the word again tomorrow
	} else {
		if card.Repetition == 0 {
			card.Interval = 1 // First correct answer, review tomorrow
		} else if card.Repetition == 1 {
			card.Interval = 6 // Second correct answer, review in 6 days
		} else {
			// Apply the spaced repetition interval, using the ease factor
			card.Interval = int64(float64(card.Interval) * card.EaseFactor)
		}
		card.Repetition++

		// Adjust the ease factor based on user feedback
		card.EaseFactor += 0.1 - float64(5-userAnswerQuality)*(0.08+float64(5-userAnswerQuality)*0.02)
		if card.EaseFactor < 1.3 {
			card.EaseFactor = 1.3 // Minimum ease factor
		}
	}

	// Calculate the next due date
	card.DueDate = time.Now().Add(time.Duration(card.Interval) * 24 * time.Hour)
}
