package oracle

// RefreshRequestInfos refreshes the info for the request with the given request id at the provided height
func (m *Module) RefreshRequestInfos(id, height int64) error {
	_, err := m.source.GetRequestStatus(height, id)
	if err != nil {
		return err
	}

	//TODO: INSERT?
	return nil
}
