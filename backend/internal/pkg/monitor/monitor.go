package monitor

type QuotaManager struct{}

func NewQuotaManager() *QuotaManager {
	return &QuotaManager{}
}

func (m *QuotaManager) CheckQuota(userID string) (bool, error) {
	return true, nil
}

func (m *QuotaManager) ConsumeQuota(userID string) error {
	return nil
}

func (m *QuotaManager) Close() error {
	return nil
}
