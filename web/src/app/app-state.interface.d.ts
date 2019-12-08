/**
 * A shape of application state
 */
interface AppState {
  /**
   * An HTTP loading indicator
   */
  isFetching: boolean;

  auth?: {
    accessToken: string,
    idToken: string,
    userInfo: UserInfo | null,
  };

  /**
   * A sidebar component storage
   */
  sidebar: {
    /**
     * A sidebar state indicator
     */
    collapsed: boolean,
  },
}
