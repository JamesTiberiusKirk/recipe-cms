package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ContextTestHelper struct {
	c   *Context
	rec *httptest.ResponseRecorder
	t   *testing.T
	req *http.Request
}

func NewContextTestHelper(t *testing.T, path string) *ContextTestHelper {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := NewCustomContext(e.NewContext(req, rec))
	c.SetPath(path)

	return &ContextTestHelper{
		c:   c,
		rec: rec,
		t:   t,
		req: req,
	}
}

func (c *ContextTestHelper) AssertAttribute(key, value string) {
	assert.Contains(c.t, c.rec.Body.String(), key+"=\""+value+"\"")
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(c.rec.Body)
	require.NoError(c.t, err)

	doc.Find("")
}
