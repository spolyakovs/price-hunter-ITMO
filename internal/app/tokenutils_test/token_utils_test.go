package tokenutils_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenutils"
)

func TestCreateDeleteTokens(t *testing.T) {
	if err := setupRedis(); err != nil {
		t.Error(err.Error())
		return
	}

	var userid uint64 = 1

	testTokenPairDetails, err := tokenutils.CreateTokens(userid)
	if err != nil {
		t.Errorf("Couldn't create tokens for userid (%d):\n\t%s", userid, err.Error())
		return
	}
	if testTokenPairDetails.AccessUuid == "" || testTokenPairDetails.RefreshUuid == "" {
		t.Errorf("Tokens for userid (%d) haven't been created:\n\t%s", userid, err.Error())
		return
	}

	if err := tokenutils.DeleteAuth(testTokenPairDetails.AccessUuid); err != nil {
		t.Errorf("Couldn't delete access token for uuid=%s:\n\t%s", testTokenPairDetails.AccessUuid, err.Error())
	}
	if err := tokenutils.DeleteAuth(testTokenPairDetails.RefreshUuid); err != nil {
		t.Errorf("Couldn't delete refresh token for uuid=%s:\n\t%s", testTokenPairDetails.RefreshUuid, err.Error())
	}
}

func TestExtractToken(t *testing.T) {
	reqEmptyAuth, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Errorf("Couldn't create for request with empty auth:\n\t%s", err.Error())
	} else {
		reqEmptyAuth.Header.Add("Authorization", "")
		tokenEmptyAuth, err := tokenutils.ExtractToken(reqEmptyAuth)
		if err == nil {
			t.Errorf("Extracted token from request with empty auth:\n\tToken: %s", tokenEmptyAuth)
		}
		if errors.Cause(err) != tokenutils.ErrTokenNotProvided {
			t.Errorf("Wrong error when extracting token from request with empty auth:\n\tToken: %s", err.Error())
		}
	}

	reqIncorrect, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Errorf("Couldn't create incorrect request:\n\t%s", err.Error())
	} else {
		reqIncorrect.Header.Add("Authorization", "Bearer")
		tokenIncorrect, err := tokenutils.ExtractToken(reqIncorrect)
		if err == nil {
			t.Errorf("Extracted token from request with incorrect auth:\n\tToken: %s", tokenIncorrect)
		}
		if errors.Cause(err) != tokenutils.ErrTokenWrongFormat {
			t.Errorf("Wrong error when extracting token from request with incorrect auth:\n\tToken: %s", err.Error())
		}
	}

	reqCorrect, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Errorf("Couldn't create correct request:\n\t%s", err.Error())
	} else {
		reqCorrect.Header.Add("Authorization", "Bearer token")
		tokenCorrect, err := tokenutils.ExtractToken(reqCorrect)
		if err != nil {
			t.Errorf("Couldn't extract token from correct request:\n\t%s", err.Error())
		}
		if tokenCorrect != "token" {
			t.Errorf("Extracted token from correct request does't match:\n\tWanted: token, Got: %s", tokenCorrect)
		}
	}
}

func TestExtractTokenMetadata(t *testing.T) {
	if err := setupRedis(); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := tokenutils.ExtractTokenMetadata("IncorrectToken"); err == nil {
		t.Errorf("Validated incorrect token (IncorrectToken)")
	} else {
		switch errors.Cause(err) {
		case tokenutils.ErrTokenDamaged:
			break
		case tokenutils.ErrTokenExpiredOrDeleted:
			t.Errorf("Incorrect token (IncorrectToken) is expired or deleted:\n\t%s", err.Error())
		default:
			t.Errorf("Something went wrong:\n\t%s", err.Error())
		}
	}

	var userid uint64 = 1

	testTokenPairDetails, err := tokenutils.CreateTokens(userid)
	if err != nil {
		t.Errorf("Couldn't create tokens for userid (%d):\n\t%s", userid, err.Error())
		return
	}
	if testTokenPairDetails.AccessToken == "" || testTokenPairDetails.AccessUuid == "" || testTokenPairDetails.RefreshUuid == "" {
		t.Errorf("Tokens for userid (%d) haven't been created:\n\t%s", userid, err.Error())
		return
	}

	if tokenDetails, err := tokenutils.ExtractTokenMetadata(testTokenPairDetails.AccessToken); err != nil {
		switch errors.Cause(err) {
		case tokenutils.ErrTokenDamaged:
			t.Errorf("Created access token (%s) for userid (%d) are not valid:\n\t%s", testTokenPairDetails.AccessToken, userid, err.Error())
		case tokenutils.ErrTokenExpiredOrDeleted:
			t.Errorf("Created access token (%s) for userid (%d) is expired (at %s) or deleted:\n\t%s", testTokenPairDetails.AccessToken, userid, time.Unix(testTokenPairDetails.AtExpires, 0), err.Error())
		default:
			t.Errorf("Something went wrong:\n\t%s", err.Error())
		}
	} else {
		if tokenDetails.Uuid != testTokenPairDetails.AccessUuid {
			t.Errorf("Got wrong tokenDetails.Uuid:\n\tWanted: %s, Got: %s", testTokenPairDetails.AccessUuid, tokenDetails.Uuid)
		}
		if tokenDetails.UserId != userid {
			t.Errorf("Got wrong tokenDetails.UserId:\n\tWanted: %d, Got: %d", userid, tokenDetails.UserId)
		}
	}

	if err := tokenutils.DeleteAuth(testTokenPairDetails.AccessUuid); err != nil {
		t.Errorf("Couldn't delete access token for uuid=%s:\n\t%s", testTokenPairDetails.AccessUuid, err.Error())
	}
	if err := tokenutils.DeleteAuth(testTokenPairDetails.RefreshUuid); err != nil {
		t.Errorf("Couldn't delete refresh token for uuid=%s:\n\t%s", testTokenPairDetails.RefreshUuid, err.Error())
	}

	if _, err := tokenutils.ExtractTokenMetadata(testTokenPairDetails.AccessToken); err == nil {
		t.Errorf("Validated deleted token  with uuid (%s)", testTokenPairDetails.AccessUuid)
	} else {
		switch errors.Cause(err) {
		case tokenutils.ErrTokenDamaged:
			t.Errorf("Created access token (%s) for userid (%d) are not valid:\n\t%s", testTokenPairDetails.AccessToken, userid, err.Error())
		case tokenutils.ErrTokenExpiredOrDeleted:
			break
		default:
			t.Errorf("Something went wrong:\n\t%s", err.Error())
		}
	}
}

func TestFetchAuth(t *testing.T) {
	if err := setupRedis(); err != nil {
		t.Error(err.Error())
	}
	tokenDetailsncorrect := &tokenutils.TokenDetails{
		Uuid:   "IncorrectUUID",
		UserId: 0,
	}

	if userIdFetched, err := tokenutils.FetchAuth(tokenDetailsncorrect); err == nil {
		t.Errorf("Fetched auth from incorrect uuid (IncorrectUUID), Got (%d)", userIdFetched)
	} else if errors.Cause(err) != tokenutils.ErrTokenExpiredOrDeleted {
		t.Errorf("Something went wrong with incorrect uuid (IncorrectUUID):\n\t%s", err.Error())
	}

	var userid uint64 = 1

	testTokenPairDetails, err := tokenutils.CreateTokens(userid)
	if err != nil {
		t.Errorf("Couldn't create tokens for userid (%d):\n\t%s", userid, err.Error())
		return
	}
	if testTokenPairDetails.AccessToken == "" || testTokenPairDetails.AccessUuid == "" || testTokenPairDetails.RefreshUuid == "" {
		t.Errorf("Tokens for userid (%d) haven't been created:\n\t%s", userid, err.Error())
		return
	}

	tokenDetailsCorrect := &tokenutils.TokenDetails{
		Uuid:   testTokenPairDetails.AccessUuid,
		UserId: userid,
	}
	userIdFetched, err := tokenutils.FetchAuth(tokenDetailsCorrect)
	if err != nil {
		t.Errorf("Couldn't fetch auth from correct uuid (%s):\n\t%s", tokenDetailsCorrect.Uuid, err.Error())
	}
	if userIdFetched != userid {
		t.Errorf("Got wrong userIdFetched:\n\tWanted: %d, Got: %d", userid, userIdFetched)
	}
}

func TestDeleteAllAuths(t *testing.T) {
	if err := setupRedis(); err != nil {
		t.Error(err.Error())
	}

	var userid uint64 = 1

	testTokenPairDetails, err := tokenutils.CreateTokens(userid)
	if err != nil {
		t.Errorf("Couldn't create tokens for userid (%d):\n\t%s", userid, err.Error())
		return
	}
	if testTokenPairDetails.AccessToken == "" || testTokenPairDetails.AccessUuid == "" || testTokenPairDetails.RefreshUuid == "" {
		t.Errorf("Tokens for userid (%d) haven't been created:\n\t%s", userid, err.Error())
		return
	}

	if err := tokenutils.DeleteAllAuths(userid); err != nil {
		t.Errorf("Couldn't delete all auths for userid (%d):\n\t%s", userid, err.Error())
	}

	tokenDetailsAccess := &tokenutils.TokenDetails{
		Uuid:   testTokenPairDetails.AccessUuid,
		UserId: userid,
	}
	if userIdFetched, err := tokenutils.FetchAuth(tokenDetailsAccess); err == nil {
		t.Errorf("Fetched deleted access auth from uuid (%s), Got (%d)", tokenDetailsAccess.Uuid, userIdFetched)
	} else if errors.Cause(err) != tokenutils.ErrTokenExpiredOrDeleted {
		t.Errorf("Something went wrong with incorrect uuid (%s):\n\t%s", tokenDetailsAccess.Uuid, err.Error())
	}

	tokenDetailsRefresh := &tokenutils.TokenDetails{
		Uuid:   testTokenPairDetails.RefreshUuid,
		UserId: userid,
	}
	if userIdFetched, err := tokenutils.FetchAuth(tokenDetailsRefresh); err == nil {
		t.Errorf("Fetched deleted refresh auth from uuid (%s), Got (%d)", tokenDetailsRefresh.Uuid, userIdFetched)
	} else if errors.Cause(err) != tokenutils.ErrTokenExpiredOrDeleted {
		t.Errorf("Something went wrong with incorrect uuid (%s):\n\t%s", tokenDetailsRefresh.Uuid, err.Error())
	}
}
