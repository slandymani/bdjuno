package oraclekeeper

/*func (k Keeper) FundOraclePool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, oracletypes.ModuleName, amount); err != nil {
		return err
	}

	oraclePool := k.GetOraclePool(ctx)
	oraclePool.DataProvidersPool = oraclePool.DataProvidersPool.Add(amount...)
	k.SetOraclePool(ctx, oraclePool)

	// store the changes in DPP
	accumulatedPaymentsForData := k.GetAccumulatedPaymentsForData(ctx)
	accumulatedPaymentsForData.AccumulatedAmount = accumulatedPaymentsForData.AccumulatedAmount.Add(amount...)

	k.SetAccumulatedPaymentsForData(ctx, accumulatedPaymentsForData)
	return nil
}

func (k Keeper) WithdrawOraclePool(ctx sdk.Context, amount sdk.Coins, recipient sdk.AccAddress) error {
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, oracletypes.ModuleName, recipient, amount); err != nil {
		return err
	}

	oraclePool := k.GetOraclePool(ctx)
	diff, hasNeg := oraclePool.DataProvidersPool.SafeSub(amount)
	if hasNeg {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "data providers pool does not have enough funds")
	}
	oraclePool.DataProvidersPool = diff
	k.SetOraclePool(ctx, oraclePool)

	// store the changes in DPP
	accumulatedPaymentsForData := k.GetAccumulatedPaymentsForData(ctx)
	accumulatedPaymentsForData.AccumulatedAmount, _ = accumulatedPaymentsForData.AccumulatedAmount.SafeSub(amount)

	k.SetAccumulatedPaymentsForData(ctx, accumulatedPaymentsForData)
	return nil
}
*/
