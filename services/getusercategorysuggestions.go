package services

import "fmt"

// GetUserCategorySuggestions returns all unique categories for the given userID.
//
// If userID is empty, an ErrInvalidUserID is returned.
// If an error occurred from the Repo, an ErrInternalRepo is returned.
func (s Service) GetUserCategorySuggestions(userID int) (suggestions []string, err error) {
	if userID < 1 {
		return nil, &ErrUserValidation{
			Field:   UserFieldID,
			Message: MessageInvalidUserID,
		}
	}

	if suggestions, err = s.repo.GetUserCategorySuggestions(userID); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return suggestions, nil
}
