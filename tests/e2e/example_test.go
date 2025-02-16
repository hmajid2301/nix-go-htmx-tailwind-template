package e2e

import (
	"testing"
)

func TestE2EExample(t *testing.T) {
	t.Run("Should do something", func(t *testing.T) {
        page, teardown := setupTest()
        t.Cleanup(func() { teardown(page) })

		expect.Locator(page.Locator("#example")).ToBeVisible()
	})
}
