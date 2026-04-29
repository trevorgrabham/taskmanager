package services

import "fmt"

// GetUserCategorySuggestions returns all unique categories for the given userID.
//
// If userID is empty, an ErrInvalidUserID is returned.
// If an error occurred from the Repo, an ErrInternalRepo is returned.
func (s Service) GetUserCategorySuggestions(userID int) (suggestions []string, err error) {
	var caller = "GetUserCategorySuggestions"
	if userID < 1 {
		return nil, fmt.Errorf("%s: %w %d", caller, ErrInvalidUserID, userID)
	}

	if suggestions, err = s.repo.GetUserCategorySuggestions(userID); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", caller, ErrInternalRepo, err)
	}

	return suggestions, nil
}
