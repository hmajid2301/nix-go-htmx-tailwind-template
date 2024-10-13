package e2e

import (
	"testing"
)

func TestE2EExample(t *testing.T) {
	t.Cleanup(ResetBrowserContexts)

	t.Run("Should do something", func(t *testing.T) {
		expect.Locator(pages[0].Locator("#example")).ToBeVisible()
	})
}
