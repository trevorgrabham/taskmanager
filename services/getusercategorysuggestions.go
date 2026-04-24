package services

import "fmt"

func (s Service) GetUserCategorySuggestions(userID int) (suggestions []string, err error) {
	var caller = "GetUserCategorySuggestions"
	if userID < 1 {
		return nil, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}

	if suggestions, err = s.repo.GetUserCategorySuggestions(userID); err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}

	return suggestions, nil
}
