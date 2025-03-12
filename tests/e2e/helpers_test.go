package e2e

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

func joinRoom(hostPlayerPage playwright.Page, otherPlayerPages []playwright.Page) error {
	// Fill host player name
	err := hostPlayerPage.GetByPlaceholder("Enter your nickname").Fill("HostPlayer")
	if err != nil {
		return fmt.Errorf("failed to fill host name: %w", err)
	}

	// Start game as host
	err = hostPlayerPage.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Start"}).Click()
	if err != nil {
		return fmt.Errorf("failed to click start button: %w", err)
	}

	// Wait for lobby to load by checking Not Ready button
	notReadyButton := hostPlayerPage.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Not Ready"})
	err = notReadyButton.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(2_000),
	})
	if err != nil {
		return fmt.Errorf("host never reached lobby: %w", err)
	}

	// Get room code from disabled input
	codeInput := hostPlayerPage.Locator("input[name='room_code']")
	code, err := codeInput.InputValue()
	if err != nil {
		return fmt.Errorf("failed to get room code: %w", err)
	}
	if code == "" {
		return fmt.Errorf("room code is empty after lobby load")
	}

	// Join other players with lobby verification
	for i, player := range otherPlayerPages {
		// Fill player details
		err := player.GetByPlaceholder("Enter your nickname").Fill(fmt.Sprintf("OtherPlayer%d", i))
		if err != nil {
			return fmt.Errorf("failed to fill player %d name: %w", i, err)
		}

		err = player.GetByPlaceholder("ABC12").Fill(code)
		if err != nil {
			return fmt.Errorf("failed to fill code for player %d: %w", i, err)
		}

		// Click join button
		err = player.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Join"}).Click()
		if err != nil {
			return fmt.Errorf("failed to click join for player %d: %w", i, err)
		}

		// Verify successful lobby entry by checking Not Ready button
		playerNotReady := player.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Not Ready"})
		err = playerNotReady.WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(2_000),
		})
		if err != nil {
			return fmt.Errorf("player %d never reached lobby: %w", i, err)
		}
	}

	return nil
}

func startGame(hostPlayerPage playwright.Page, otherPlayerPages []playwright.Page) error {
	err := joinRoom(hostPlayerPage, otherPlayerPages)
	if err != nil {
		return err
	}

	for _, player := range append(otherPlayerPages, hostPlayerPage) {
		err = player.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Ready"}).Click()
		if err != nil {
			return err
		}
	}

	err = hostPlayerPage.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Start Game"}).Click()
	if err != nil {
		return err
	}

	return nil
}

func getPlayerRoles(
	hostPlayerPage playwright.Page,
	otherPlayerPages []playwright.Page,
) (playwright.Page, []playwright.Page, error) {
	var fibber playwright.Page
	normals := make([]playwright.Page, 0)
	allPlayers := append([]playwright.Page{hostPlayerPage}, otherPlayerPages...)

	// Wait for all players to reach the answer submission screen
	for _, player := range allPlayers {
		err := player.Locator("#submit_answer_form").WaitFor(
			playwright.LocatorWaitForOptions{
				Timeout: playwright.Float(15_000), // 15 seconds timeout
			},
		)
		if err != nil {
			return nil, nil, fmt.Errorf("player didn't reach answer screen: %w", err)
		}
	}

	// Wait for role assignment to be complete
	submitButton := hostPlayerPage.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Not Ready"})
	err := submitButton.WaitFor(
		playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(10_000), // 10 seconds timeout
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("submit button never appeared: %w", err)
	}

	// Find fibber with proper waiting
	for _, player := range allPlayers {
		fibberLocator := player.GetByText("You are fibber")

		// Wait for fibber text with shorter per-player timeout
		err := fibberLocator.WaitFor(
			playwright.LocatorWaitForOptions{
				Timeout: playwright.Float(1_000), // 5 seconds per player
			},
		)

		if err == nil {
			fibber = player
			break
		}
	}

	if fibber == nil {
		return nil, nil, fmt.Errorf("fibber not found after checking all players")
	}

	// Collect normal players
	for _, player := range allPlayers {
		if player != fibber {
			normals = append(normals, player)
		}
	}

	// Verify we found exactly 1 fibber
	if len(normals) != len(allPlayers)-1 {
		return nil, nil, fmt.Errorf("expected %d normal players but found %d",
			len(allPlayers)-1, len(normals))
	}

	return fibber, normals, nil
}
