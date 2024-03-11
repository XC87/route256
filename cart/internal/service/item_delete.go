package service

func (cartService *CartService) DeleteItem(userId int64, skuId int64) error {
	if userId <= 0 {
		return ErrUserInvalid
	}

	err := cartService.repository.DeleteItem(userId, skuId)
	if err != nil {
		return err
	}
	return nil
}
