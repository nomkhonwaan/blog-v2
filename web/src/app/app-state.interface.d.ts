/**
 * A shape of the application state
 */
interface AppState {
  /**
   * An HTTP loading indicator
   */
  isFetching: boolean,

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
