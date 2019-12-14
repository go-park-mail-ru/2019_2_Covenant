package repository

import (
	"2019_2_Covenant/internal/subscriptions"
	. "2019_2_Covenant/tools/vars"
	"database/sql"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) subscriptions.Repository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (ssR *SubscriptionRepository) Subscribe(userID uint64, subscriptionID uint64) error {
	if err := ssR.db.QueryRow("SELECT id FROM subscriptions WHERE user_id = $1 AND subscribed_to = $2",
		userID,
		subscriptionID,
	).Scan(); err == nil {
		return ErrAlreadyExist
	}

	if _, err := ssR.db.Exec("INSERT INTO subscriptions (user_id, subscribed_to) VALUES ($1, $2)",
		userID,
		subscriptionID,
	); err != nil {
			return err
	}

	return nil
}

func (ssR *SubscriptionRepository) Unsubscribe(userID uint64, subscriptionID uint64) error {
	res, err := ssR.db.Exec("DELETE FROM subscriptions WHERE user_id = $1 AND subscribed_to = $2",
		userID,
		subscriptionID,
	)

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}

	return nil
}
