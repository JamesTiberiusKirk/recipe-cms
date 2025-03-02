package playground

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: this is just a playground for figuring out testing strat
func Test_testPage(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := common.NewCustomContext(e.NewContext(req, rec))
	c.SetPath("/pg/testpage")

	route := TestRoute{}
	err := route.HandleTestPage(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, c.Response().Status)
	require.Equal(t, echo.MIMETextHTML, c.Response().Header().Get(echo.HeaderContentType))

	doc, err := goquery.NewDocumentFromReader(rec.Body)
	require.NoError(t, err)

	form := doc.Find("#recipe_form")
	assert.Equal(t, 1, form.Length(), "need to have recipe_form id")

	ta := form.Find("#test_area-textarea")
	assert.Equal(t, 1, ta.Length(), "need to have test_area id inside the form")
}
