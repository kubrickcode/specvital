package api

// MarshalJSON implements json.Marshaler for GetSpecDocumentByRepository200JSONResponse.
// This is needed because the generated type alias doesn't inherit the MarshalJSON method
// from RepoSpecDocumentResponse, causing the unexported union field to be ignored.
func (r GetSpecDocumentByRepository200JSONResponse) MarshalJSON() ([]byte, error) {
	return RepoSpecDocumentResponse(r).MarshalJSON()
}
