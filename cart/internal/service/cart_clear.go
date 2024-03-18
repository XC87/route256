package service

func (cartService *CartService) DeleteItemsByUserId(userId int64) error {
	if userId <= 0 {
		return ErrUserInvalid
	}

	err := cartService.repository.DeleteItemsByUserId(userId)
	if err != nil {
		return err
	}
	return nil
}
