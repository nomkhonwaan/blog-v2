/**
 * GitHub Gist in JSON format
 */
interface Gist {
  /**
   * A description
   */
  description: string

  /**
   * Visibility indicator
   */
  public: boolean

  /**
   * List of files in the Gist
   */
  files: string[]

  /**
   * To-be rendered HTML elements
   */
  div: string

  /**
   * Link to the CSS stylesheet
   */
  stylesheet: string
}
