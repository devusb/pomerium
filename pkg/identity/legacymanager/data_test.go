package legacymanager

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pomerium/pomerium/pkg/grpc/session"
	"github.com/pomerium/pomerium/pkg/protoutil"
)

func TestUser_UnmarshalJSON(t *testing.T) {
	var u User
	err := json.Unmarshal([]byte(`{
		"name": "joe",
		"email": "joe@test.com",
		"some-other-claim": "xyz"
	}`), &u)
	assert.NoError(t, err)
	assert.NotNil(t, u.User)
	assert.Equal(t, "joe", u.User.Name)
	assert.Equal(t, "joe@test.com", u.User.Email)
	assert.Equal(t, map[string]*structpb.ListValue{
		"some-other-claim": {Values: []*structpb.Value{protoutil.ToStruct("xyz")}},
	}, u.Claims)
}

func TestSession_NextRefresh(t *testing.T) {
	tm1 := time.Date(2020, 6, 5, 12, 0, 0, 0, time.UTC)
	s := Session{
		Session:         &session.Session{},
		lastRefresh:     tm1,
		gracePeriod:     time.Second * 10,
		coolOffDuration: time.Minute,
	}
	assert.Equal(t, tm1.Add(time.Minute), s.NextRefresh())

	tm2 := time.Date(2020, 6, 5, 13, 0, 0, 0, time.UTC)
	s.OauthToken = &session.OAuthToken{
		ExpiresAt: timestamppb.New(tm2),
	}
	assert.Equal(t, tm2.Add(-time.Second*10), s.NextRefresh())

	tm3 := time.Date(2020, 6, 5, 12, 15, 0, 0, time.UTC)
	s.ExpiresAt = timestamppb.New(tm3)
	assert.Equal(t, tm3, s.NextRefresh())
}

func TestSession_UnmarshalJSON(t *testing.T) {
	tm := time.Date(2020, 6, 5, 12, 0, 0, 0, time.UTC)
	var s Session
	err := json.Unmarshal([]byte(`{
		"iss": "https://some.issuer.com",
		"sub": "subject",
		"exp": `+fmt.Sprint(tm.Unix())+`,
		"iat": `+fmt.Sprint(tm.Unix())+`,
		"some-other-claim": "xyz"
	}`), &s)
	assert.NoError(t, err)
	assert.NotNil(t, s.Session)
	assert.NotNil(t, s.Session.IdToken)
	assert.Equal(t, "https://some.issuer.com", s.Session.IdToken.Issuer)
	assert.Equal(t, "subject", s.Session.IdToken.Subject)
	assert.Equal(t, timestamppb.New(tm), s.Session.IdToken.ExpiresAt)
	assert.Equal(t, timestamppb.New(tm), s.Session.IdToken.IssuedAt)
	assert.Equal(t, map[string]*structpb.ListValue{
		"some-other-claim": {Values: []*structpb.Value{protoutil.ToStruct("xyz")}},
	}, s.Claims)
}
