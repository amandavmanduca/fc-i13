package entity

type Investor struct {
	ID string
	Name string
	AssetPosition []*InvestorAssetPosition // endereço na memória - ponteiro
}

type InvestorAssetPosition struct {
	AssetID string
	Shares int
}

func NewInvestor(id string) *Investor {
	return &Investor {
		ID: id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

// método pra uma struct - função normal mas pra ser acessada tem que passar um i?
func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	i.AssetPosition = append(i.AssetPosition, assetPosition) // adicionando o asset no "array"
}

func (i *Investor) UpdateAssetPosition(assetID string, shareQuant int) {
	assetPosition := i.GetAssetPosition(assetID);
	if assetPosition == nil {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetID, shareQuant))
	} else {
		assetPosition.Shares += shareQuant
	}
}

func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assetID {
			return assetPosition
		}
	}
	return nil
}

func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition {
		AssetID: assetID,
		Shares: shares,
	}
}
