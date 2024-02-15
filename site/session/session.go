package session

import (
	"fmt"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

const (
	sessionName                  = "session"
	sessionLifeSpan              = 216000 // 1 hour
	shortcodeLoginSessionDefault = "unauthed"
)

// Manager maintains a record of open sessions.
type Manager struct {
	Jar *sessions.CookieStore
}

var (
	// NOTE: temporary store, ideally we need something like redis or
	// memcache which would also have a TTL on the record
	// tbh having this in mem wodnt even be too bad, we just need to cleanup
	// TODO: either way ill need to wrap this in an interface and dependency wrapper
	// so then it can be swapped out for whatever store mechanism we want
	shortcodeLoginSessions = map[string]string{}
	shortCodeChannels      = map[string](chan string){}
)

// New returns an instantiated session manager.
func New() *Manager {
	authKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)

	return &Manager{
		Jar: sessions.NewCookieStore(authKey, encryptionKey),
	}
}

// InitSession will store a new session or refresh an existing one.
func (m *Manager) InitSession(username string, c echo.Context) {
	sess, _ := m.Jar.Get(c.Request(), sessionName)

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionLifeSpan,
		HttpOnly: true,
	}

	sess.Values["username"] = username
	_ = sess.Save(c.Request(), c.Response())
}

// InitShortCodeSess initialises a short code session on for the original device.
func (m *Manager) InitShortCodeSess(c echo.Context) string {
	short := shortid.MustGenerate()

	sess, _ := m.Jar.Get(c.Request(), sessionName)

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15, // 15 minutes
		HttpOnly: true,
	}

	sess.Values["shortid"] = short

	shortcodeLoginSessions[short] = shortcodeLoginSessionDefault
	shortCodeChannels[short] = make(chan string)

	_ = sess.Save(c.Request(), c.Response())

	return short
}

// AuthShortCodeSess runs for the second device with the session already ongoing.
// to run for second device already logged in.
func (m *Manager) AuthShortCodeSess(short string, c echo.Context) error {
	// BUG: this always triggers
	// if m.IsAuthenticated(c, false) {
	// 	return fmt.Errorf("not authenticated")
	// }

	usernameInShort, ok := shortcodeLoginSessions[short]
	if !ok {
		return fmt.Errorf("no short code session")
	}

	if usernameInShort != shortcodeLoginSessionDefault {
		return fmt.Errorf("shortcode session already claimed")
	}

	fmt.Printf("short code sess map %+v", shortcodeLoginSessions)

	username, err := m.GetUser(c)
	if err != nil {
		logrus.Errorf("error getting username %s", err.Error())
		return fmt.Errorf("error getting user %w", err)
	}

	fmt.Printf("short code sess map %+v", shortcodeLoginSessions)

	shortcodeLoginSessions[short] = username
	shortCodeChannels[short] <- "authenticated"
	return nil
}

func (m *Manager) GetShortCodeChan(c echo.Context) (<-chan string, error) {
	sess, err := m.Jar.Get(c.Request(), sessionName)
	if err != nil {
		logrus.Errorf("error getting session %s", err.Error())
		return nil, fmt.Errorf("error getting session %w", err)
	}

	shortI, ok := sess.Values["shortid"]
	if !ok {
		logrus.Errorf("could not get short")
		return nil, fmt.Errorf("error get short")
	}

	short, ok := shortI.(string)
	if !ok {
		logrus.Errorf("could not type cast")
		return nil, fmt.Errorf("error type casting")
	}

	ch, ok := shortCodeChannels[short]
	if !ok || ch == nil {
		return nil, fmt.Errorf("could not find channel")
	}

	return ch, nil
}

// TerminateSession will cease tracking the session for the current user.
func (m *Manager) TerminateSession(c echo.Context) {
	sess, _ := m.Jar.Get(c.Request(), sessionName)
	// MaxAge < 0 means delete imediately
	sess.Options.MaxAge = -1
	_ = sess.Save(c.Request(), c.Response())
}

// IsAuthenticated checks that a provided request is born from an active session.
// As long as there is an active session, true is returned, else false.
func (m *Manager) IsAuthenticated(c echo.Context, checkAndLoginShortSession bool) bool {
	sess, err := m.Jar.Get(c.Request(), sessionName)
	if err != nil {
		logrus.Errorf("error getting session %s", err.Error())
		return false
	}

	if sess.Values["username"] == nil && checkAndLoginShortSession {
		if sess.Values["shortid"] == nil {
			return false
		}

		short, ok := sess.Values["shortid"].(string)
		if !ok {
			return false
		}

		usernameInShort, ok := shortcodeLoginSessions[short]
		if !ok {
			return false
		}

		fmt.Printf("short %s \n", short)
		fmt.Printf("usernameInShort %s \n", usernameInShort)

		if usernameInShort == "" {
			return false
		}

		if usernameInShort == shortcodeLoginSessionDefault {
			return false
		}

		m.InitSession(usernameInShort, c)

		// map cleanup
		delete(shortcodeLoginSessions, short)
		delete(shortCodeChannels, short)

		return true
	}

	if sess.Values["username"] == nil {
		fmt.Printf("username is nil\n")
		return false
	}
	return true
}

// GetUser checks that a provided request is born from an active session.
// As long as there is an active session, User is returned, else empty User.
func (m *Manager) GetUser(c echo.Context) (string, error) {
	sess, err := m.Jar.Get(c.Request(), sessionName)
	if err != nil {
		return "", err
	}

	if sess.Values == nil {
		return "", nil
	}

	usernameI, ok := sess.Values["username"]
	if !ok {
		return "", nil
	}

	username, ok := usernameI.(string)
	if !ok {
		return "", nil
	}

	return username, nil
}
