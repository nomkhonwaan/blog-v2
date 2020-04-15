/**
 * An authentication result which includes OAuth2 and OpenID Connect properties
 */
interface AuthResult {
  /**
   * The access_token issued by the authorization server
   */
  accessToken: string;

  /**
   * ID Tokens are used in token-based authentication to cache user profile information
   * and provide it to a client application
   */
  idToken: string;

  /**
   * A special kind of token that is used to authenticate a user without them needing to re-authenticate
   */
  refreshToken: string;

  /**
   * If the "state" parameter was present in the client authorization request.
   * The exact value received from the client.
   */
  state: string;

  /**
   * The lifetime in seconds of the access token.
   * For example, the value "3600" denotes that the access token will
   * expire in one hour from the time the response was generated.
   * If omitted, the authorization server SHOULD provide the
   * expiration time via other means or document the default value.
   */
  expiresIn: number;

  /**
   * The type of the token issued as described in
   */
  tokenType: string;

  /**
   * If identical to the scope requested by the client; otherwise, REQUIRED
   */
  scope: string;
}
