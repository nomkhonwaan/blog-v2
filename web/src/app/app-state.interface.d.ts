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
   * A content component
   */
  content: {
    sidebar: {
      collapsed: boolean,
    },
  };

  /**
   * An admin component
   */
  admin: {
    editor: {
      sidebar: {
        collapsed: boolean,
      },
    },
  };
}
