package openid

import (
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

type jwtTokenValidatorMock struct {
	t     *testing.T
	Calls chan Call
}

func newJwtTokenValidatorMock(t *testing.T) *jwtTokenValidatorMock {
	return &jwtTokenValidatorMock{t, make(chan Call)}
}

type validateCall struct {
	req *http.Request
	t   string
}

type validateResp struct {
	jt  *jwt.Token
	err error
}

func (j *jwtTokenValidatorMock) validate(r *http.Request, t string) (*jwt.Token, error) {
	j.Calls <- &validateCall{r, t}
	vr := (<-j.Calls).(*validateResp)
	return vr.jt, vr.err
}

func (j *jwtTokenValidatorMock) assertValidate(t string, jt *jwt.Token, err error) {
	call := (<-j.Calls).(*validateCall)
	if t != anything && call.t != t {
		j.t.Error("Expected validate with token", t, "but was", call.t)
	}
	j.Calls <- &validateResp{jt, err}
}

func (j *jwtTokenValidatorMock) close() {
	close(j.Calls)
}

func (j *jwtTokenValidatorMock) assertDone() {
	if _, more := <-j.Calls; more {
		j.t.Fatal("Did not expect more calls.")
	}
}
