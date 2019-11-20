/**
 * A group of posts regarded as having particular shared characteristics
 */
interface Category {
  /**
   * Name of the category
   */
  name: string

  /**
   * Valid URL string composes with name and ID
   */
  slug: string

  /**
   * List of latest published posts are belongging to the category
   */
  latestPublishedPosts: Post[];
}
