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
   * A sidebar component
   */
  sidebar: {
    /**
     * A sidebar showing or hiding state
     */
    collapsed: boolean,
  };

  /**
   * An editor component
   */
  editor: {
    /**
     * An editor's sidebar component
     */
    sidebar: {
      /**
       * A sidebar showing or hiding state
       */
      collapsed: boolean,
    },
  };
}
