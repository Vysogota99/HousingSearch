package jwt

// TokenDetails house these tokens definitions, their expiration periods and uuids
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AccTExpires  int64
	RefTExpires  int64
}

// AccessDetails ...
type AccessDetails struct {
	AccessUUID      string
	UserID          int64
	TelephoneNumber string
}
