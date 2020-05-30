package pgkit_test

import (
	"testing"

	"github.com/JordanOcokoljic/pgkit"
)

func TestParseDetails(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		user     string
		password string
		location string
		port     string
		database string
	}{
		{
			name:     "NoValues",
			uri:      "postgresql://",
			user:     "",
			password: "",
			location: "",
			port:     "",
			database: "",
		},
		{
			name:     "Host",
			uri:      "postgresql://localhost",
			user:     "",
			password: "",
			location: "localhost",
			port:     "",
			database: "",
		},
		{
			name:     "HostAndPort",
			uri:      "postgresql://localhost:5432",
			user:     "",
			password: "",
			location: "localhost",
			port:     "5432",
			database: "",
		},
		{
			name:     "HostAndDatabase",
			uri:      "postgresql://localhost/mydb",
			user:     "",
			password: "",
			location: "localhost",
			port:     "",
			database: "mydb",
		},
		{
			name:     "UserAndHost",
			uri:      "postgresql://user@localhost",
			user:     "user",
			password: "",
			location: "localhost",
			port:     "",
			database: "",
		},
		{
			name:     "UserPasswordAndHost",
			uri:      "postgresql://user:secret@localhost",
			user:     "user",
			password: "secret",
			location: "localhost",
			port:     "",
			database: "",
		},
		{
			name:     "UserHostAndDatabase",
			uri:      "postgresql://other@localhost/otherdb",
			user:     "other",
			password: "",
			location: "localhost",
			port:     "",
			database: "otherdb",
		},
		{
			name:     "UserPasswordHostAndDatabase",
			uri:      "postgresql://other:password@localhost/otherdb",
			user:     "other",
			password: "password",
			location: "localhost",
			port:     "",
			database: "otherdb",
		},
		{
			name:     "All",
			uri:      "postgresql://user:secret@localhost:5432/postgres",
			user:     "user",
			password: "secret",
			location: "localhost",
			port:     "5432",
			database: "postgres",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(s *testing.T) {
			cd, err := pgkit.ParseDetails(test.uri)
			if err != nil {
				s.Fatalf(err.Error())
			}

			if cd.User != test.user {
				s.Errorf(
					"user not extracted correctly: was %s expected %s",
					cd.User, test.user,
				)
			}

			if cd.Password != test.password {
				s.Errorf(
					"password not extracted correctly: was %s expected %s",
					cd.Password, test.password,
				)
			}

			if cd.Location != test.location {
				s.Errorf(
					"location not extracted correctly: was %s expected %s",
					cd.Location, test.location,
				)
			}

			if cd.Port != test.port {
				s.Errorf(
					"port not extracted correctly: was %s expected %s",
					cd.Port, test.port,
				)
			}

			if cd.Database != test.database {
				s.Errorf(
					"database not extracted correctly: was %s expected %s",
					cd.Database, test.database,
				)
			}
		})
	}
}

func TestUriParseDetailsOptions(t *testing.T) {
	con := "postgresql://localhost?sslmode=disabled&application_name=pgkit"
	cd, err := pgkit.ParseDetails(con)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(cd.Options) != 2 {
		t.Fatalf("expected 2 options but there were %d", len(cd.Options))
	}

	if cd.Options["sslmode"] != "disabled" {
		t.Errorf(
			"sslmode not extracted correctly, was %s expected disabled",
			cd.Options["sslmode"],
		)
	}

	if cd.Options["application_name"] != "pgkit" {
		t.Errorf(
			"application_name not extracted correctly, was %s expected pgkit",
			cd.Options["application_name"],
		)
	}
}

func TestConnectionDetailCopy(t *testing.T) {
	cd := pgkit.NewConnectionDetail()
	cd.User = "user"
	cd.Password = "password"
	cd.Location = "location"
	cd.Port = "port"
	cd.Database = "database"
	cd.Options["option"] = "value"

	ncd := cd.Copy()
	ncd.User = "user2"
	ncd.Password = "password2"
	ncd.Location = "location2"
	ncd.Port = "port2"
	ncd.Database = "database2"
	ncd.Options["option"] = "value2"

	if cd.User == ncd.User {
		t.Errorf("user field changed on original instance")
	}

	if cd.Password == ncd.Password {
		t.Errorf("password field changed on original instance")
	}

	if cd.Location == ncd.Location {
		t.Errorf("location field changed on original instance")
	}

	if cd.Port == ncd.Port {
		t.Errorf("port field changed on original instance")
	}

	if cd.Database == ncd.Database {
		t.Errorf("database field changed on original instance")
	}

	if cd.Options["option"] == ncd.Options["option"] {
		t.Errorf("option field changed on original instance")
	}
}

func TestConnectionIsValid(t *testing.T) {
	cd := pgkit.NewConnectionDetail()
	cd.Password = "password"
	if cd.IsValid() {
		t.Errorf("detail was valid with unset username and set password")
	}

	cd = pgkit.NewConnectionDetail()
	cd.Port = "5432"
	if cd.IsValid() {
		t.Errorf("detail was valid with unset location and set port")
	}

	cd = pgkit.NewConnectionDetail()
	cd.Port = "port"
	if cd.IsValid() {
		t.Errorf("detail was valid with non numerical port")
	}

	cd = pgkit.NewConnectionDetail()
	cd.Options[""] = "value"
	if cd.IsValid() {
		t.Errorf("detail was valid with blank option key")
	}

	cd = pgkit.NewConnectionDetail()
	cd.User = "user"
	cd.Password = "password"
	cd.Location = "location"
	cd.Port = "5432"
	cd.Database = "database"
	cd.Options["sslmode"] = "disable"

	if !cd.IsValid() {
		t.Errorf("detail was invalid when all values set correctly")
	}
}
