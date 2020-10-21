package mongohelper

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func newClient(host, username, password, source string, minPoolSize, maxPoolSize uint64, connIdeaTime time.Duration, readMode readpref.Mode) (*mongo.Client, error) {
	ops := options.Client().SetHosts([]string{host}).SetConnectTimeout(20 * time.Second)

	if username != "" {
		cred := options.Credential{}
		cred.Username = username
		cred.Password = password
		cred.AuthMechanism = "SCRAM-SHA-1"
		cred.AuthSource = source

		ops.SetAuth(cred)
	}

	if maxPoolSize > 0 {
		ops.SetMaxPoolSize(maxPoolSize)
	}
	if minPoolSize > 0 {
		ops.SetMinPoolSize(minPoolSize)
	}

	if connIdeaTime > 0 {
		ops.SetMaxConnIdleTime(connIdeaTime)
	}

	pref, err := readpref.New(readMode)
	if err != nil {
		return nil, err
	}
	ops.SetReadPreference(pref)

	client, err := mongo.NewClient(ops)
	if err != nil {
		return nil, err
	}

	if err := client.Connect(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

func NewSecondaryClient(host, username, password, source string) (*mongo.Client, error) {
	return newClient(host, username, password, source, 16, 64, time.Minute, readpref.SecondaryPreferredMode)
}

func NewClient(host, username, password, source string) (*mongo.Client, error) {
	return newClient(host, username, password, source, 16, 64, time.Minute, readpref.PrimaryMode)
}
